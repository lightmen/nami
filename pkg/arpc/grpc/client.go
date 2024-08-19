package grpc

import (
	"context"
	"time"

	"github.com/lightmen/nami/alog"
	"github.com/lightmen/nami/codec"
	"github.com/lightmen/nami/message"
	"github.com/lightmen/nami/metadata"
	"github.com/lightmen/nami/pkg/arpc"
	"github.com/lightmen/nami/pkg/endpoint"
	"github.com/lightmen/nami/registry"
	"github.com/lightmen/nami/transport/agrpc/balancer"
)

func Init(dis registry.Discovery) {
	arpc.SetDiscorey(dis)
	cli := New(WithDiscovery(dis))

	arpc.SetDefaultClient(cli)
}

type client struct {
	dis registry.Discovery
}

func New(opts ...clientOption) arpc.Client {
	cli := &client{
		dis: arpc.GetDiscorey(),
	}

	for _, opt := range opts {
		opt(cli)
	}

	return cli
}

func (cli *client) Request(ctx context.Context, target, route, uid string, cmd int32, req, rsp codec.Codec, opts ...arpc.CallOption) (err error) {
	info := &arpc.CallInfo{
		Target: target,
		Route:  route,
		Cmd:    cmd,
		Req:    req,
		Rsp:    rsp,
		UID:    uid,
		Type:   message.REQUEST,
	}

	return cli.call(ctx, info, opts...)
}

func (cli *client) Event(ctx context.Context, target, route, uid string, cmd int32, req codec.Codec, opts ...arpc.CallOption) (err error) {
	info := &arpc.CallInfo{
		Target: target,
		Route:  route,
		Cmd:    cmd,
		Req:    req,
		UID:    uid,
		Type:   message.EVENT,
	}
	return cli.call(ctx, info, opts...)
}

func (cli *client) Notify(ctx context.Context, target string, uid string, cmd int32, req codec.Codec, opts ...arpc.CallOption) (err error) {
	info := &arpc.CallInfo{
		Target: target,
		Route:  uid,
		Cmd:    cmd,
		Req:    req,
		Type:   message.NOTIFY,
		UID:    uid,
	}

	//如果只有一个玩家

	return cli.call(ctx, info, opts...)
}

func (cli *client) NotifyAll(ctx context.Context, srv string, cmd int32, req codec.Codec, opts ...arpc.CallOption) (err error) {
	info := &arpc.CallInfo{
		Route: "",
		Cmd:   cmd,
		Req:   req,
		Type:  message.NOTIFYALL,
	}
	instances, err := cli.dis.GetService(ctx, srv)
	if err != nil {
		return
	}

	for _, instance := range instances {
		info.Target = endpoint.GetGrpcEndpoint(instance.Endpoints)

		err = cli.call(ctx, info, opts...)
		if err != nil {
			return
		}
	}

	return
}

func (cli *client) Broadcast(ctx context.Context, srv, uid string, cmd int32, req codec.Codec, opts ...arpc.CallOption) (err error) {
	info := &arpc.CallInfo{
		Route: uid,
		UID:   uid,
		Cmd:   cmd,
		Req:   req,
		Type:  message.EVENT,
	}
	instances, err := cli.dis.GetService(ctx, srv)
	if err != nil {
		return
	}

	for _, instance := range instances {
		info.Target = endpoint.GetGrpcEndpoint(instance.Endpoints)

		err = cli.call(ctx, info, opts...)
		if err != nil {
			return
		}
	}

	return
}

func (cli *client) call(ctx context.Context, info *arpc.CallInfo, opts ...arpc.CallOption) (err error) {
	//初始化变量
	rsp := info.Rsp
	target := info.Target
	for _, opt := range opts {
		opt(info)
	}
	uid := info.UID

	//构建packet包
	mpkt, err := cli.buildPacket(info)
	if err != nil {
		alog.ErrorCtx(ctx, "%s|%d|%s|buildPacket error: %s", uid, info.Cmd, target, err.Error())
		return
	}

	//构建context
	ctx, err = cli.buildContext(ctx, info)
	if err != nil {
		alog.ErrorCtx(ctx, "%s|%d|%s|buildContext error: %s", uid, info.Cmd, target, err.Error())
		return
	}

	//获取链接
	conn, err := GetConn(ctx, target, info.Addr)
	if err != nil {
		alog.ErrorCtx(ctx, "%s|%d|%s|GetConn error: %s", uid, info.Cmd, target, err.Error())
		return err
	}
	msgClient := message.NewMessageClient(conn)

	reply, err := msgClient.HandleMessage(ctx, mpkt)
	if err != nil {
		alog.ErrorCtx(ctx, "%s|%d|%s|HandleMessage error: %s", uid, info.Cmd, target, err.Error())
		return
	}
	if mpkt.Head.Type == message.REQUEST {
		err = rsp.Unmarshal(reply.Body)
		if err != nil {
			alog.ErrorCtx(ctx, "%s|%d|%s|Unmarshal error: %s", uid, info.Cmd, target, err.Error())
			return
		}
	}

	return
}

func (cli *client) buildContext(ctx context.Context, info *arpc.CallInfo) (context.Context, error) {
	uid := info.UID
	if info.Addr == "" { //说明不是直连
		ctx = balancer.NewNameContext(ctx, balancer.Consistent)
		ctx = balancer.NewTargetContext(ctx, info.Target)
		ctx = balancer.NewParamContext(ctx, info.Route)
	}

	if metadata.GetUIDFromClientContext(ctx) != uid {
		ctx = metadata.AppendUIDToClientContext(ctx, uid)
	}

	return ctx, nil
}

func (cli *client) buildPacket(param *arpc.CallInfo) (*message.Packet, error) {
	reqBuf, err := param.Req.Marshal()
	if err != nil {
		return nil, err
	}

	mpkt := &message.Packet{
		Head: &message.Head{
			Seq:     time.Now().UnixMilli(),
			Route:   param.Route,
			Cmd:     param.Cmd,
			Type:    param.Type,
			Targets: []string{param.UID},
		},
		Body: reqBuf,
	}

	return mpkt, nil
}
