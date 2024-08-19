package rcontext

import (
	"context"

	"google.golang.org/grpc/metadata"
)

const (
	headKey = "pbrpc_head-bin"
	uidKey  = "pbrpc_uid"
)

// NewContext 存储uid和PacketHead到context中
func NewContext(parent context.Context, uid string) (ctx context.Context) {

	// ctx = metadata.AppendToOutgoingContext(parent, headKey, string(buf))
	ctx = metadata.AppendToOutgoingContext(parent, uidKey, uid)

	return
}

func GetUID(ctx context.Context) (string, bool) {
	vals := metadata.ValueFromIncomingContext(ctx, uidKey)
	if len(vals) > 0 {
		return vals[0], true
	}

	return "", false
}
