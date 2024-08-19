package agrpc

import (
	"context"
	"fmt"
	"time"

	"github.com/lightmen/nami/middleware"
	"github.com/lightmen/nami/middleware/metadata"
	"github.com/lightmen/nami/middleware/tracing"
	"github.com/lightmen/nami/registry"
	"github.com/lightmen/nami/transport"
	"github.com/lightmen/nami/transport/agrpc/resolver/discovery"
	"google.golang.org/grpc"
	grpcinsecure "google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	grpcmd "google.golang.org/grpc/metadata"
)

type ClientOption func(o *clientOptions)

type clientOptions struct {
	endpoint     string
	timeout      time.Duration
	discovery    registry.Discovery
	middleware   []middleware.Middleware
	ints         []grpc.UnaryClientInterceptor
	grpcOpts     []grpc.DialOption
	balancerName string
}

func WithEndpoint(endpoint string) ClientOption {
	return func(o *clientOptions) {
		o.endpoint = endpoint
	}
}

func WithUnaryInterceptor(ints ...grpc.UnaryClientInterceptor) ClientOption {
	return func(o *clientOptions) {
		o.ints = ints
	}
}

func WithDiscovery(discovery registry.Discovery) ClientOption {
	return func(o *clientOptions) {
		o.discovery = discovery
	}
}

func WithOptions(opts ...grpc.DialOption) ClientOption {
	return func(o *clientOptions) {
		o.grpcOpts = opts
	}
}

func WithBalancerName(name string) ClientOption {
	return func(o *clientOptions) {
		o.balancerName = name
	}
}

// WithMiddleware with client middleware.
func WithMiddleware(m ...middleware.Middleware) ClientOption {
	return func(o *clientOptions) {
		o.middleware = m
	}
}

func defaultClientOptions() clientOptions {
	var kacp = keepalive.ClientParameters{
		Time:                10 * time.Second, // send pings every 10 seconds if there is no activity
		Timeout:             2 * time.Second,  // wait 2 second for ping ack before considering the connection dead
		PermitWithoutStream: true,             // send pings even without active streams
	}

	options := clientOptions{
		timeout: 2500 * time.Millisecond,
		middleware: []middleware.Middleware{
			tracing.Client(),
			metadata.Client(),
		},
		grpcOpts: []grpc.DialOption{
			grpc.WithKeepaliveParams(kacp),
		},
	}

	return options
}

// Dial returns a GRPC connection.
func Dial(ctx context.Context, opts ...ClientOption) (*grpc.ClientConn, error) {
	return dial(ctx, opts...)
}

func dial(ctx context.Context, opts ...ClientOption) (*grpc.ClientConn, error) {
	options := defaultClientOptions()

	for _, o := range opts {
		o(&options)
	}

	ints := []grpc.UnaryClientInterceptor{
		unaryClientInterceptor(options.middleware, options.timeout),
	}
	if len(options.ints) > 0 {
		ints = append(ints, options.ints...)
	}
	grpcOpts := []grpc.DialOption{
		grpc.WithChainUnaryInterceptor(ints...),
		grpc.WithTransportCredentials(grpcinsecure.NewCredentials()),
	}

	if options.balancerName != "" {
		grpcOpts = append(grpcOpts,
			grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"loadBalancingConfig": [{"%s":{}}]}`, options.balancerName)),
		)
	}

	if options.discovery != nil {
		grpcOpts = append(grpcOpts,
			grpc.WithResolvers(
				discovery.NewBuilder(
					options.discovery,
				)))
	}

	if len(options.grpcOpts) > 0 {
		grpcOpts = append(grpcOpts, options.grpcOpts...)
	}

	return grpc.DialContext(ctx, options.endpoint, grpcOpts...)
}

func unaryClientInterceptor(ms []middleware.Middleware, timeout time.Duration) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		op := GetCmd(ctx, nil)
		if op == "" {
			op = method
		}
		ctx = transport.NewClientContext(ctx, &Transport{
			operation: op,
			endpoint:  cc.Target(),
			header:    headerCarrier{},
		})

		if timeout > 0 {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, timeout)
			defer cancel()
		}

		h := func(ctx context.Context, req interface{}) (interface{}, error) {
			if tr, ok := transport.FromClientContext(ctx); ok {
				header := tr.Header()
				keys := header.Keys()
				keyvals := make([]string, 0, len(keys))
				for _, k := range keys {
					keyvals = append(keyvals, k, header.Get(k))
				}
				ctx = grpcmd.AppendToOutgoingContext(ctx, keyvals...)
			}
			return reply, invoker(ctx, method, req, reply, cc, opts...)
		}
		if len(ms) > 0 {
			h = middleware.Chain(ms...)(h)
		}

		_, err := h(ctx, req)
		return err
	}
}
