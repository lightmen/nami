package chain

import "net/http"

type Middleware func(http.Handler) http.Handler

type Chain interface {
	Append(midlewares ...Middleware) Chain
	Prepend(midlewares ...Middleware) Chain
	Then(h http.Handler) http.Handler
	ThenFunc(fn http.HandlerFunc) http.Handler
}

type chain struct {
	middlewares []Middleware
}

func New(middlewares ...Middleware) Chain {
	return &chain{middlewares: append([]Middleware{}, middlewares...)}
}

func (c *chain) Append(middlewares ...Middleware) Chain {
	c.middlewares = join(c.middlewares, middlewares)
	return c
}

func (c *chain) Prepend(middlewares ...Middleware) Chain {
	c.middlewares = join(middlewares, c.middlewares)
	return c
}

func (c *chain) Then(h http.Handler) http.Handler {
	if h == nil {
		h = http.DefaultServeMux
	}

	for i := range c.middlewares {
		h = c.middlewares[len(c.middlewares)-1-i](h)
	}

	return h
}

func (c *chain) ThenFunc(fn http.HandlerFunc) http.Handler {
	if fn == nil {
		return c.Then(nil)
	}

	return c.Then(fn)
}

func join(a, b []Middleware) []Middleware {
	c := make([]Middleware, 0, len(a)+len(b))
	c = append(c, a...)
	c = append(c, b...)
	return c
}
