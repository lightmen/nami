package recovery

import (
	"context"
	"runtime/debug"

	"github.com/lightmen/nami/alog"
	"github.com/lightmen/nami/middleware"
)

func Recovery() middleware.Middleware {
	return func(next middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (any, error) {
			defer func() {
				if result := recover(); result != nil {
					alog.ErrorCtx(ctx, "%v\n%s", result, debug.Stack())
					return
				}
			}()

			return next(ctx, req)
		}
	}
}
