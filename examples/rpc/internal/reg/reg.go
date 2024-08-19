package reg

import (
	"time"

	"github.com/lightmen/nami/alog"
	"github.com/lightmen/nami/pkg/arpc"
	"github.com/lightmen/nami/registry"
	"github.com/lightmen/nami/registry/etcd"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func GetRegistry(addr string) (reg registry.Registrar, err error) {
	if addr == "" {
		addr = "127.0.0.1:2379"
	}
	reg, err = CreateEtcdReg(addr)

	dis, ok := reg.(registry.Discovery)
	if ok {
		arpc.SetDiscorey(dis)
	}

	return
}

func CreateEtcdReg(addr string) (reg registry.Registrar, err error) {
	cfg := clientv3.Config{
		Endpoints:            []string{addr},
		DialTimeout:          time.Second * 10,
		DialKeepAliveTime:    time.Second * 10,
		DialKeepAliveTimeout: time.Second * 10,
	}

	etcdCli, err := clientv3.New(cfg)
	if err != nil {
		alog.Error("clientv3.New error: %s", err.Error())
		return
	}

	reg = etcd.New(etcdCli)
	return
}
