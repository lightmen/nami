package acontext

import (
	"context"
	"sync"
	"time"
)

var (
	_ context.Context = (*Context)(nil)
)

type Context struct {
	ctx     context.Context
	keyvals sync.Map
}

func New(opt ...Option) *Context {
	ctx := &Context{
		ctx: context.Background(),
	}

	for _, o := range opt {
		o(ctx)
	}

	return ctx
}

func WithValue(parent context.Context, key, val any) context.Context {
	if ctx, ok := parent.(*Context); ok {
		return ctx.WithValue(key, val)
	}

	return context.WithValue(parent, key, val)
}

func (c *Context) WithValue(key, val any) context.Context {
	c.keyvals.Store(key, val)
	return c
}

func (c *Context) Value(key any) any {
	if val, ok := c.keyvals.Load(key); ok {
		return val
	}

	return c.ctx.Value(key)
}

func (c *Context) Deadline() (deadline time.Time, ok bool) {
	return c.ctx.Deadline()
}

func (c *Context) Done() <-chan struct{} {
	return c.ctx.Done()
}

func (c *Context) Err() error {
	return c.ctx.Err()
}
