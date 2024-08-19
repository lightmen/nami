package acontext

import (
	"context"
	"sync"
)

type keyvalKey struct{}
type carrier struct {
	keyvals sync.Map
}

func NewCarrierContext(ctx context.Context) context.Context {
	val := &carrier{}
	return context.WithValue(ctx, keyvalKey{}, val)
}

func GetCarrierFromContext(ctx context.Context) (val *carrier, ok bool) {
	val, ok = ctx.Value(keyvalKey{}).(*carrier)
	if !ok {
		return
	}

	return
}

func DataFromCarrierCtx(ctx context.Context, key any) (val any, ok bool) {
	carrier, ok := GetCarrierFromContext(ctx)
	if !ok {
		return nil, false
	}

	return carrier.keyvals.Load(key)
}

func DataToCarrierCtx(ctx context.Context, key, val any) context.Context {
	carrier, ok := GetCarrierFromContext(ctx)
	if !ok {
		return context.WithValue(ctx, key, val)
	}

	carrier.keyvals.Store(key, val)
	return ctx
}
