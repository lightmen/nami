package arpc

import (
	"github.com/lightmen/nami/codec"
	"github.com/lightmen/nami/message"
)

type CallInfo struct {
	Target string      //必填，目标服务名,当Addr不为空时，请求根据一致性hash发送到某个服务
	UID    string      //必填，表示发送消息的玩家id，可以为空
	Route  string      //必填，用于一致性hash路由，可以为空
	Cmd    int32       //必填，请求命令字
	Req    codec.Codec //必填，请求数据
	Rsp    codec.Codec //选填，回包
	Addr   string      //选填，服务地址，当该地址不为空的时候，请求发送到该地址所在的服务， 其地址格式类似：grpc://127.0.0.1:32521
	Type   message.Type
}

type CallOption func(o *CallInfo)

func WithAddr(addr string) CallOption {
	return func(o *CallInfo) {
		o.Addr = addr
	}
}
