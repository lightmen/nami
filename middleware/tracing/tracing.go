package tracing

import (
	"context"

	"github.com/lightmen/nami/middleware"
	"github.com/lightmen/nami/transport"
	"go.opentelemetry.io/otel/trace"
)

func Server(opt ...Option) middleware.Middleware {
	tracer := NewTracer(trace.SpanKindServer, opt...)

	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (rsp any, err error) {
			if tr, ok := transport.FromServerContext(ctx); ok {
				var span trace.Span
				ctx, span = tracer.Start(ctx, tr.Operation(), tr.Header())
				setServerSpan(ctx, span)

				defer func() {
					tracer.End(ctx, span, rsp, err)
				}()
			}
			return handler(ctx, req)
		}
	}
}

func Client(opt ...Option) middleware.Middleware {
	tracer := NewTracer(trace.SpanKindClient, opt...)

	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (rsp any, err error) {
			if tr, ok := transport.FromClientContext(ctx); ok {
				var span trace.Span
				ctx, span = tracer.Start(ctx, tr.Operation(), tr.Header())
				defer func() {
					tracer.End(ctx, span, rsp, err)
				}()
			}

			return handler(ctx, req)
		}
	}
}
