package discovery

import (
	"context"
	"errors"
	"time"

	"github.com/lightmen/nami/alog"
	"github.com/lightmen/nami/pkg/endpoint"
	"github.com/lightmen/nami/registry"
	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/resolver"
)

type discoveryResolver struct {
	watchName string
	w         registry.Watcher
	cc        resolver.ClientConn

	ctx    context.Context
	cancel context.CancelFunc
}

func (r *discoveryResolver) watch() {
	for {
		select {
		case <-r.ctx.Done():
			return
		default:
		}
		ins, err := r.w.Next()
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return
			}
			alog.Error("[resolver][%s] Failed to watch discovery endpoint: %v", r.watchName, err)
			time.Sleep(time.Second)
			continue
		}
		r.update(ins)
	}
}

func (r *discoveryResolver) update(ins []*registry.Instance) {
	addrs := make([]resolver.Address, 0)
	endpoints := make(map[string]struct{})
	for _, in := range ins {
		ept, err := endpoint.ParseEndpoint(in.Endpoints, "grpc")
		if err != nil {
			alog.Error("[resolver][%s] Failed to parse discovery endpoint: %v", r.watchName, err)
			continue
		}
		if ept == "" {
			continue
		}
		// filter redundant endpoints
		if _, ok := endpoints[ept]; ok {
			continue
		}

		endpoints[ept] = struct{}{}
		addr := resolver.Address{
			ServerName: in.Name,
			Attributes: parseAttributes(in.MetaData),
			Addr:       ept,
		}
		addr.Attributes = addr.Attributes.WithValue("rawServiceInstance", in)
		addrs = append(addrs, addr)
	}
	if len(addrs) == 0 {
		alog.Info("[resolver][%s] Zero endpoint found,refused to write, instances: %v", r.watchName, ins)
		return
	}

	err := r.cc.UpdateState(resolver.State{Addresses: addrs})
	if err != nil {
		alog.Error("[resolver][%s] failed to update state: %s", r.watchName, err)
		return
	}

	// alog.Debug("[resolver][%s] update addrs: %d, ins: %s", r.watchName, len(addrs), cast.ToJson(ins))
}

func (r *discoveryResolver) Close() {
	r.cancel()
	err := r.w.Stop()
	if err != nil {
		alog.Error("[resolver] failed to watch top: %s", err)
	}
}

// ResolveNow 空函数，仅仅是为了实现 resolver.Resolver 接口
func (r *discoveryResolver) ResolveNow(resolver.ResolveNowOptions) {

}

func parseAttributes(md map[string]string) *attributes.Attributes {
	var a *attributes.Attributes
	for k, v := range md {
		if a == nil {
			a = attributes.New(k, v)
		} else {
			a = a.WithValue(k, v)
		}
	}
	return a
}
