package metactx

import (
	"context"

	"github.com/lightmen/nami/middleware"
)

// Interceptor 用于设置context.Context
type Interceptor func(ctx context.Context) (context.Context, error)

// Wrapper 将context进行包装
func Wrapper(ints ...Interceptor) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (reply any, err error) {
			for _, inter := range ints {
				ctx, err = inter(ctx)
				if err != nil {
					return
				}
			}
			return handler(ctx, req)
		}
	}
}
