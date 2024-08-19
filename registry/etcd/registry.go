package etcd

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/lightmen/nami/alog"
	"github.com/lightmen/nami/registry"
)

var (
	_ registry.Registrar = (*Registry)(nil)
	_ registry.Discovery = (*Registry)(nil)
)

// Option is etcd registry option.
type Option func(o *options)

type options struct {
	ctx       context.Context
	namespace string
	ttl       time.Duration
	maxRetry  int
}

// Context with registry context.
func Context(ctx context.Context) Option {
	return func(o *options) {
		o.ctx = ctx
	}
}

// Namespace with registry namespace.
func Namespace(ns string) Option {
	return func(o *options) {
		o.namespace = ns
	}
}

// RegisterTTL with register ttl.
func RegisterTTL(ttl time.Duration) Option {
	return func(o *options) {
		o.ttl = ttl
	}
}

func MaxRetry(num int) Option {
	return func(o *options) {
		o.maxRetry = num
	}
}

// Registry is etcd registry.
type Registry struct {
	opts   *options
	client *clientv3.Client
	kv     clientv3.KV
	lease  clientv3.Lease
}

// New creates etcd registry
func New(client *clientv3.Client, opts ...Option) (r *Registry) {
	op := &options{
		ctx:       context.Background(),
		namespace: "/servers/node",
		ttl:       time.Second * 15,
		maxRetry:  5,
	}
	for _, o := range opts {
		o(op)
	}
	return &Registry{
		opts:   op,
		client: client,
		kv:     clientv3.NewKV(client),
	}
}

// Register the registration.
func (r *Registry) Register(ctx context.Context, service *registry.Instance) error {
	key := fmt.Sprintf("%s/%s/%s", r.opts.namespace, service.Name, service.ID)
	value, err := marshal(service)
	if err != nil {
		return err
	}
	if r.lease != nil {
		r.lease.Close()
	}
	r.lease = clientv3.NewLease(r.client)
	leaseID, err := r.registerWithKV(ctx, key, value)
	if err != nil {
		return err
	}

	go r.heartBeat(r.opts.ctx, leaseID, key, value)
	return nil
}

// Unregister the registration.
func (r *Registry) Unregister(ctx context.Context, service *registry.Instance) error {
	defer func() {
		if r.lease != nil {
			r.lease.Close()
		}
	}()
	key := fmt.Sprintf("%s/%s/%s", r.opts.namespace, service.Name, service.ID)

	alog.InfoCtx(ctx, "[etcd]unregister key: %s", key)

	_, err := r.client.Delete(ctx, key)
	return err
}

// GetService return the service instances in memory according to the service name.
func (r *Registry) GetService(ctx context.Context, name string) ([]*registry.Instance, error) {
	key := fmt.Sprintf("%s/%s", r.opts.namespace, name)
	resp, err := r.kv.Get(ctx, key, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	items := make([]*registry.Instance, 0, len(resp.Kvs))
	for _, kv := range resp.Kvs {
		si, err := unmarshal(kv.Value)
		if err != nil {
			return nil, err
		}
		if name != "" && si.Name != name {
			continue
		}
		items = append(items, si)
	}
	return items, nil
}

// Watch creates a watcher according to the service name.
func (r *Registry) Watch(ctx context.Context, name string) (registry.Watcher, error) {
	key := fmt.Sprintf("%s/%s", r.opts.namespace, name)
	return newWatcher(ctx, key, name, r.client)
}

// registerWithKV create a new lease, return current leaseID
func (r *Registry) registerWithKV(ctx context.Context, key string, value string) (clientv3.LeaseID, error) {
	grant, err := r.lease.Grant(ctx, int64(r.opts.ttl.Seconds()))
	if err != nil {
		return 0, err
	}
	_, err = r.client.Put(ctx, key, value, clientv3.WithLease(grant.ID))
	if err != nil {
		return 0, err
	}
	return grant.ID, nil
}

func (r *Registry) heartBeat(ctx context.Context, leaseID clientv3.LeaseID, key string, value string) {
	curLeaseID := leaseID
	kac, err := r.client.KeepAlive(ctx, leaseID)
	if err != nil {
		alog.ErrorCtx(ctx, "%s|keep alive error: %s", key, err.Error())
		curLeaseID = 0
	}

	for {
		if curLeaseID == 0 {
			curLeaseID, kac = r.tryRegister(ctx, key, value)
			if kac == nil {
				alog.ErrorCtx(ctx, "%s|tryRegister failed, exit", key)
				return
			}
		}

		select {
		case _, ok := <-kac:
			if !ok {
				if ctx.Err() != nil {
					// channel closed due to context cancel
					return
				}
				// need to retry registration
				curLeaseID = 0
				continue
			}
		case <-ctx.Done():
			return
		}
	}
}

func (r *Registry) tryRegister(ctx context.Context, key, value string) (curLeaseID clientv3.LeaseID, kac <-chan *clientv3.LeaseKeepAliveResponse) {
	var err error
	var retreat []int

	maxSleepTime := int(r.opts.ttl.Seconds())
	maxSleepTime = maxSleepTime / 2
	if maxSleepTime <= 0 {
		maxSleepTime = 1
	}

	// try to registerWithKV
	for retryCnt := 0; retryCnt < r.opts.maxRetry; retryCnt++ {
		if ctx.Err() != nil {
			return 0, nil
		}
		// prevent infinite blocking
		idChan := make(chan clientv3.LeaseID, 1)
		errChan := make(chan error, 1)
		cancelCtx, cancel := context.WithCancel(ctx)
		go func() {
			defer cancel()
			id, registerErr := r.registerWithKV(cancelCtx, key, value)
			if registerErr != nil {
				errChan <- registerErr
			} else {
				idChan <- id
			}
		}()

		select {
		case <-time.After(3 * time.Second):
			cancel()
			continue
		case <-errChan:
			continue
		case curLeaseID = <-idChan:
		}

		kac, err = r.client.KeepAlive(ctx, curLeaseID)
		if err == nil {
			break
		}

		if maxSleepTime >= 1<<retryCnt {
			retreat = append(retreat, 1<<retryCnt)
		}

		sleepIdx := rand.Intn(len(retreat))

		alog.ErrorCtx(ctx, "%s|%d|keep alive error: %s, sleep seconds: %d", key, retryCnt, err.Error(), retreat[sleepIdx])

		time.Sleep(time.Duration(retreat[sleepIdx]) * time.Second)
	}

	if _, ok := <-kac; !ok {
		// retry failed
		alog.Debug("%s|keep alive failed , exit heartbeat", key)
		return 0, nil
	}

	return
}
