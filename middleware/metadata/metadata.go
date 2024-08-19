package metadata

import (
	"context"

	"github.com/lightmen/nami/metadata"
	"github.com/lightmen/nami/middleware"
	"github.com/lightmen/nami/transport"
)

// Server server-side,将transporter的header里面的key-value写到metadata里
func Server(opts ...Option) middleware.Middleware {
	opt := &options{
		prefix: []string{"x-md"},
	}
	for _, o := range opts {
		o(opt)
	}

	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (rsp any, err error) {
			tr, ok := transport.FromServerContext(ctx)
			if !ok {
				return handler(ctx, req)
			}

			md := opt.md.Clone()

			header := tr.Header()
			for _, key := range header.Keys() {
				if opt.hasPrefix(key) {
					md.Set(key, header.Get(key))
				}
			}

			ctx = metadata.NewServerContext(ctx, md)

			ctx = metadata.AddUIDFromMDCtx(ctx)

			return handler(ctx, req)
		}
	}
}

// Client client-side, 将metadata的key-value写到transporter的header里面去
func Client(opts ...Option) middleware.Middleware {
	opt := &options{
		prefix: []string{"x-md-global"},
	}
	for _, o := range opts {
		o(opt)
	}

	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (rsp any, err error) {
			tr, ok := transport.FromClientContext(ctx)
			if !ok {
				return handler(ctx, req)
			}

			header := tr.Header()

			for key, val := range opt.md {
				header.Set(key, val)
			}

			// x-md-global-
			header = opt.setServerHeader(ctx, header)

			//因为client中的md与server中的md的key有出现重复的可能，重新key，我们采用client的。
			//因此，这里我们把client的md的set放在server后面，是为了在出现重复key的情况下，client的md
			//可以覆盖server的md
			header = opt.setClientHeader(ctx, header)

			ctx = metadata.AddUIDFromMDCtx(ctx)
			return handler(ctx, req)
		}
	}
}

func (opt *options) setServerHeader(ctx context.Context, header transport.Header) transport.Header {
	if md, ok := metadata.FromServerContext(ctx); ok {
		for k, v := range md {
			if opt.hasPrefix(k) {
				header.Set(k, v)
			}
		}
	}

	return header
}

func (opt *options) setClientHeader(ctx context.Context, header transport.Header) transport.Header {
	if md, ok := metadata.FromClientContext(ctx); ok {
		for key, val := range md {
			header.Set(key, val)
		}
	}

	return header
}
