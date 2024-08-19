package grpc

import (
	"context"
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/lightmen/nami/alog"
	"github.com/lightmen/nami/pkg/arpc"
	"github.com/lightmen/nami/pkg/endpoint"
	"github.com/lightmen/nami/pkg/safe"
	"github.com/lightmen/nami/registry"
	"github.com/lightmen/nami/transport/agrpc"
	"github.com/lightmen/nami/transport/agrpc/balancer"
	"github.com/lightmen/nami/transport/agrpc/resolver/discovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

var (
	poolMap  sync.Map
	poolLock sync.Mutex
)

const (
	maxPoolExistTime = time.Second * 360

	// KeepAliveTime is the duration of time after which if the client doesn't see
	// any activity it pings the server to see if the transport is still alive.
	KeepAliveTime = time.Duration(10) * time.Second

	// KeepAliveTimeout is the duration of time for which the client waits after having
	// pinged for keepalive check and if no activity is seen even after that the connection
	// is closed.
	KeepAliveTimeout = time.Duration(3) * time.Second
)

type poolWrap struct {
	// pool.Pool     //@deprecated
	conn          *grpc.ClientConn
	lastExistTime time.Time
}

func init() {
	safe.Go(poolWork)
}

func poolWork() {
	time.Sleep(time.Second * 5) //这里先睡眠5秒，等待 arpc.GetDiscorey()返回的值不为空

	dis := arpc.GetDiscorey()

	//该for循环做两个事：
	//1. 从服务发现里面，获取所有服务的Instance实例集合，并且根据实例集合更新pool集合里面pool的存在时间
	//2. 变量pool集合，如果发现pool超过 maxPoolExistTime 时间没有更新存在时间的，删掉这个pool
	for {
		time.Sleep(time.Second * 30)

		instances, err := dis.GetService(context.Background(), "")
		if err != nil {
			alog.Error("dis.GetService error: %s", err.Error())
			return
		}

		nowTime := time.Now()
		for _, instance := range instances {
			updatePoolExist(instance, nowTime)
		}

		delConns := make([]*grpc.ClientConn, 0)
		poolMap.Range(func(key, value any) bool {
			wrap := value.(*poolWrap)
			//超过 maxPoolExistTime 时间没有被服务发现获取到，默认为这个pool不存在，可以删掉了
			if nowTime.Sub(wrap.lastExistTime) >= maxPoolExistTime {
				delConns = append(delConns, wrap.conn)
				poolMap.Delete(key)
				alog.Info("delete grpc client pool: %v", key)
			}
			return true
		})

		for _, conn := range delConns {
			conn.Close()
		}
	}
}

// updatePoolExist 更新poolMap的存在时间
func updatePoolExist(instance *registry.Instance, nowTime time.Time) {
	var pool *poolWrap

	//一个instance实例可能在poolMap中以两种方式存在：其一是服务名，其二是grpc地址
	targets := []string{instance.Name, endpoint.GetGrpcEndpoint(instance.Endpoints)}
	for _, target := range targets {
		val, ok := poolMap.Load(target)
		if ok {
			pool = val.(*poolWrap)
			pool.lastExistTime = nowTime
		}
	}
}

// 获取连接， target可以为srv名字或者grpc地址，类似：gamesrv 或者 grpc://192.168.15.117:33308
func getClientConn(ctx context.Context, target string) (wrap *poolWrap, err error) {
	val, ok := poolMap.Load(target)
	if ok {
		return val.(*poolWrap), nil
	}

	poolLock.Lock()
	defer poolLock.Unlock()

	val, ok = poolMap.Load(target)
	if ok {
		return val.(*poolWrap), nil
	}

	conn, err := createClientConn(ctx, target)
	if err != nil {
		return nil, err
	}

	wrap = &poolWrap{
		conn:          conn,
		lastExistTime: time.Now(),
	}

	poolMap.Store(target, wrap)

	return wrap, nil
}

func createClientConn(ctx context.Context, target string) (conn *grpc.ClientConn, err error) {
	u, err := url.Parse(target)
	if err != nil {
		return nil, err
	}

	opts := []agrpc.ClientOption{
		agrpc.WithOptions(
			grpc.WithKeepaliveParams(keepalive.ClientParameters{
				Time:                KeepAliveTime,
				Timeout:             KeepAliveTimeout,
				PermitWithoutStream: true,
			}),
		),
	}

	if u.Scheme == "grpc" {
		opts = append(opts,
			agrpc.WithEndpoint(u.Host),
		)
	} else {
		opts = append(opts,
			agrpc.WithEndpoint(fmt.Sprintf("%s:///%s", discovery.Name, target)),
			agrpc.WithDiscovery(arpc.GetDiscorey()),
			agrpc.WithBalancerName(balancer.Mix),
		)
	}

	return agrpc.Dial(ctx, opts...)
}

// // GetPool 根据target获取pool
// func GetPool(target string) (pool.Pool, error) {
// 	val, ok := poolMap.Load(target)
// 	if ok {
// 		return val.(*poolWrap), nil
// 	}

// 	poolLock.Lock()
// 	defer poolLock.Unlock()

// 	val, ok = poolMap.Load(target)
// 	if ok {
// 		return val.(*poolWrap), nil
// 	}

// 	pool, err := createPool(target)
// 	if err != nil {
// 		return nil, err
// 	}

// 	wrap := &poolWrap{
// 		Pool:          pool,
// 		lastExistTime: time.Now(),
// 	}
// 	poolMap.Store(target, wrap)

// 	return wrap, nil
// }

// func createPool(target string) (pool.Pool, error) {
// 	u, err := url.Parse(target)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if u.Scheme == "grpc" {
// 		return pool.New(u.Host)
// 	}

// 	addr := fmt.Sprintf("%s:///%s", discovery.Name, target)
// 	return pool.New(addr, pool.WithDial(dialForServer))
// }

// func dialForServer(address string) (*grpc.ClientConn, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	return agrpc.Dial(
// 		ctx,
// 		agrpc.WithEndpoint(address),
// 		agrpc.WithDiscovery(arpc.GetDiscorey()),
// 	)
// }
