package etcd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/lightmen/nami/registry"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func (r *Registry) Register(ctx context.Context, instance *registry.Instance) (err error) {
	content, err := json.Marshal(instance)
	if err != nil {
		return
	}

	//设置租约时间
	grant, err := r.client.Grant(context.Background(), int64(r.opts.ttl.Seconds()))
	if err != nil {
		return err
	}

	key := fmt.Sprintf("%s/%s/%s", r.opts.namespace, instance.Name, instance.ID)
	//注册服务并绑定租约
	_, err = r.client.Put(r.opts.ctx, key, string(content), clientv3.WithLease(grant.ID))
	if err != nil {
		return
	}

	//设置续租 定期发送请求
	leaseRespChan, err := r.client.KeepAlive(r.opts.ctx, grant.ID)
	if err != nil {
		return err
	}
	r.leaseID = grant.ID
	r.keepAliveChan = leaseRespChan

	go r.listenLease()

	return
}

func (r *Registry) listenLease() {
	for {
		select {
		case <-r.opts.ctx.Done():
			return
		case rsp, ok := <-r.keepAliveChan:
			if !ok {
				r.opts.log.Info("lease end, rsp: %s", rsp.String())
				return
			}
			r.opts.log.Debug("continue lease: %s", rsp.String())
		}
	}
}

func (r *Registry) Unregister(ctx context.Context, instance *registry.Instance) (err error) {
	_, err = r.client.Revoke(r.opts.ctx, r.leaseID)
	if err != nil {
		return
	}

	return r.client.Close()
}
