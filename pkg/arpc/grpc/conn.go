package grpc

import (
	"context"

	"google.golang.org/grpc"
)

// GetConn 根据target获取conn，target为srv名字, 类似 gamesrv, addr为地址,类似 grpc://192.168.15.117:33308，addr可以为空
func GetConn(ctx context.Context, target, addr string) (*grpc.ClientConn, error) {
	if addr != "" {
		target = addr
	}
	wrap, err := getClientConn(ctx, target)
	if err != nil {
		return nil, err
	}

	return wrap.conn, nil
}
