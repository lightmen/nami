package ratelimit

import (
	"context"

	"github.com/lightmen/nami/middleware"
	"go.uber.org/ratelimit"
)

type Option func(o *options)

func Limiter(limiter ratelimit.Limiter) Option {
	return func(o *options) {
		o.limiter = limiter
	}
}

type options struct {
	limiter ratelimit.Limiter
}

func Server(opts ...Option) middleware.Middleware {
	opt := &options{
		limiter: ratelimit.New(10000),
	}

	for _, o := range opts {
		o(opt)
	}

	return func(next middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (any, error) {
			opt.limiter.Take()

			return next(ctx, req)
		}
	}
}
