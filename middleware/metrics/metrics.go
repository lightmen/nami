package metrics

import (
	"context"
	"time"

	"github.com/lightmen/nami/metrics"
	"github.com/lightmen/nami/middleware"
	"github.com/lightmen/nami/pkg/aerror"
	"github.com/lightmen/nami/pkg/cast"
	"github.com/lightmen/nami/transport"
)

type Option func(*option)

func WithRequest(requests metrics.Counter) Option {
	return func(o *option) {
		o.requests = requests
	}
}

func WithSeconds(seconds metrics.Observer) Option {
	return func(o *option) {
		o.seconds = seconds
	}
}

type option struct {
	// counter: <client/server>_cmd_requests_total{cmd, code}
	requests metrics.Counter
	// histogram: <client/server>_cmd_durations_bucket{cmd}
	seconds metrics.Observer
}

func Server(opts ...Option) middleware.Middleware {
	o := &option{}
	for _, opt := range opts {
		opt(o)
	}

	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (any, error) {
			var (
				cmd  string
				code int32
			)
			startTime := time.Now()
			if info, ok := transport.FromServerContext(ctx); ok {
				cmd = info.Operation()
			}
			reply, err := handler(ctx, req)
			if err != nil {
				code = aerror.Code(err)
			}

			if o.requests != nil {
				o.requests.With(cmd, cast.ToString(code)).Inc()
			}

			if o.seconds != nil {
				o.seconds.With(cmd).Observe(float64(time.Since(startTime).Milliseconds()))
			}

			return reply, err
		}
	}
}

func Client(opts ...Option) middleware.Middleware {
	o := &option{}
	for _, opt := range opts {
		opt(o)
	}

	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (rsp any, err error) {
			var (
				cmd  string
				code int32
			)
			startTime := time.Now()
			if info, ok := transport.FromServerContext(ctx); ok {
				cmd = info.Operation()
			}
			reply, err := handler(ctx, req)
			if err != nil {
				code = aerror.Code(err)
			}

			if o.requests != nil {
				o.requests.With(cmd, cast.ToString(code)).Inc()
			}

			if o.seconds != nil {
				o.seconds.With(cmd).Observe(float64(time.Since(startTime).Milliseconds()))
			}

			return reply, err
		}
	}
}
