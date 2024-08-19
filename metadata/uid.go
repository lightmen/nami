package metadata

import (
	"context"

	"github.com/lightmen/nami/pkg/acontext"
)

type uidKey struct{}

func NewUIDContext(ctx context.Context, uid string) context.Context {
	ctx = acontext.WithValue(ctx, uidKey{}, uid)
	return ctx
}

// FromUIDContext 从uidKey结构里面获取uid
func FromUIDContext(ctx context.Context) (uid string, ok bool) {
	uid, ok = ctx.Value(uidKey{}).(string)
	return
}

// GetUID 从先从uidKey结构获取uid，如果uid不存在，再从metadata里面获取uid
func GetUID(ctx context.Context) string {
	if uid, ok := FromUIDContext(ctx); ok {
		return uid
	}

	return GetUIDFromMDCtx(ctx)
}

// AddUIDFromMDCtx 如果uidKey结构的ctx中没有uid，就从metadata中提取uid，存储到uidKey中
func AddUIDFromMDCtx(ctx context.Context) context.Context {
	if _, ok := FromUIDContext(ctx); ok {
		return ctx
	}

	if uid := GetUIDFromMDCtx(ctx); uid != "" {
		ctx = NewUIDContext(ctx, uid)
	}

	return ctx
}

// GetUIDFromMDCtx 根据ctx获取metadata，再从metadata里面获取uid
func GetUIDFromMDCtx(ctx context.Context) string {
	//优先调用 FromServerContext， 因为大多数情况下，GetUIDFromMDCtx是作为
	//server端被调用
	if md, ok := FromServerContext(ctx); ok {
		if uid := md.Get(UIDKey); uid != "" {
			return uid
		}
	}

	return GetUIDFromClientContext(ctx)
}

func GetUIDFromClientContext(ctx context.Context) string {
	if md, ok := FromClientContext(ctx); ok {
		if uid := md.Get(UIDKey); uid != "" {
			return uid
		}
	}

	return ""
}

func AppendUIDToClientContext(ctx context.Context, uid string) context.Context {
	return AppendToClientContext(ctx, UIDKey, uid)
}

func NewClientUIDContext(ctx context.Context, uid string) context.Context {
	return AppendToClientContext(ctx, UIDKey, uid)
}
