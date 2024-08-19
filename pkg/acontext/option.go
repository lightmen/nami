package acontext

import "context"

type Option func(*Context)

func WithCtx(ctx context.Context) Option {
	return func(c *Context) {
		c.ctx = ctx
	}
}
