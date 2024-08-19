package arpc

import (
	"context"

	"github.com/lightmen/nami/codec"
)

type Client interface {
	// Request 发送 req 到target服务， srv为服务名, 如果要发送到特定地址的服务，需要带上 WithAddr(addr)这个CallOption
	Request(ctx context.Context, srv, route, uid string, cmd int32, req, rsp codec.Codec, opts ...CallOption) (err error)
	// Request 发送 event 到target服务， srv为服务名
	Event(ctx context.Context, srv, route, uid string, cmd int32, event codec.Codec, opts ...CallOption) (err error)
	// Broadcast 将 req 广播到所有的srv服务上
	Broadcast(ctx context.Context, srv, uid string, cmd int32, req codec.Codec, opts ...CallOption) (err error)
	//Notify 对某几个玩家发送通知
	Notify(ctx context.Context, srv string, uid string, cmd int32, req codec.Codec, opts ...CallOption) (err error)
	//NotifyAll 对所有在线玩家发送通知
	NotifyAll(ctx context.Context, srv string, cmd int32, req codec.Codec, opts ...CallOption) (err error)
}

var defaultClient Client

func SetDefaultClient(cli Client) {
	defaultClient = cli
}

// Request 发送 req 到target服务， srv为服务名，如果要发送到特定地址的服务，需要带上 WithAddr(addr)这个CallOption
func Request(ctx context.Context, srv, route, uid string, cmd int32, req, rsp codec.Codec, opts ...CallOption) (err error) {
	return defaultClient.Request(ctx, srv, route, uid, cmd, req, rsp, opts...)
}

// Request 发送 event 到target服务， srv为服务名
func Event(ctx context.Context, srv, route, uid string, cmd int32, event codec.Codec, opts ...CallOption) (err error) {
	return defaultClient.Event(ctx, srv, route, uid, cmd, event, opts...)
}

// Notify 对某几个玩家发送通知
func Notify(ctx context.Context, target string, uid string, cmd int32, req codec.Codec, opts ...CallOption) (err error) {
	return defaultClient.Notify(ctx, target, uid, cmd, req, opts...)
}

// NotifyAll 对所有在线玩家发送通知
func NotifyAll(ctx context.Context, srv string, cmd int32, req codec.Codec, opts ...CallOption) (err error) {
	return defaultClient.NotifyAll(ctx, srv, cmd, req, opts...)
}

// Broadcast 将 req 广播到所有的srv服务上
func Broadcast(ctx context.Context, srv, uid string, cmd int32, req codec.Codec, opts ...CallOption) (err error) {
	return defaultClient.Broadcast(ctx, srv, uid, cmd, req, opts...)
}
