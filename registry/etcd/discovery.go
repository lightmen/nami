package etcd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/lightmen/nami/registry"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func (r *Registry) Watch(ctx context.Context, srvName string) (err error) {
	prefix := fmt.Sprintf("%s/%s", r.opts.namespace, srvName)

	//根据前缀获取现有的key
	resp, err := r.client.Get(context.Background(), prefix, clientv3.WithPrefix())
	if err != nil {
		return err
	}

	for _, ev := range resp.Kvs {
		r.setService(string(ev.Key), ev.Value)
	}

	//监视前缀，修改变更的server
	go r.watcher(prefix)

	return
}

func (r *Registry) GetService(ctx context.Context, srvName string) (instances []*registry.Instance, err error) {
	r.lk.RLock()
	instances = make([]*registry.Instance, 0, len(r.instances))
	for _, inst := range r.instances {
		instances = append(instances, inst)
	}
	r.lk.RUnlock()

	return
}

//watcher
func (r *Registry) watcher(prefix string) {
	rch := r.client.Watch(context.Background(), prefix, clientv3.WithPrefix())
	r.opts.log.Info("watching prefix:%s now...", prefix)
	for wresp := range rch {
		for _, ev := range wresp.Events {
			switch ev.Type {
			case mvccpb.PUT: //修改或者新增
				r.setService(string(ev.Kv.Key), ev.Kv.Value)
			case mvccpb.DELETE: //删除
				r.delService(string(ev.Kv.Key))
			}
		}
	}
}

func (r *Registry) setService(key string, buf []byte) {
	r.lk.Lock()
	defer r.lk.Unlock()

	inst := &registry.Instance{}
	err := json.Unmarshal(buf, inst)
	if err != nil {
		return
	}

	if r.instances == nil {
		r.instances = make(map[string]*registry.Instance)
	}

	r.instances[key] = inst
}

func (r *Registry) delService(key string) {
	r.lk.Lock()
	defer r.lk.Unlock()

	if r.instances == nil {
		return
	}

	delete(r.instances, key)
}
