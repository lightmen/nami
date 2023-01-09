package etcd

import (
	"context"
	"sync"
	"time"

	"github.com/lightmen/nami/core/log"
	"github.com/lightmen/nami/registry"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var (
	_ registry.Registrar = (*Registry)(nil)
	_ registry.Discovery = (*Registry)(nil)
)

type Registry struct {
	opts   *options
	client *clientv3.Client
	//define members related to service register
	leaseID       clientv3.LeaseID
	keepAliveChan <-chan *clientv3.LeaseKeepAliveResponse
	//define members related to service discovery
	instances map[string]*registry.Instance
	lk        sync.RWMutex
}

func New(client *clientv3.Client, opts ...Option) *Registry {
	op := &options{
		ctx:       context.Background(),
		namespace: "/server/node",
		ttl:       time.Second * 5,
		log:       log.Default(),
	}
	for _, o := range opts {
		o(op)
	}

	r := &Registry{
		opts:   op,
		client: client,
	}

	return r
}
