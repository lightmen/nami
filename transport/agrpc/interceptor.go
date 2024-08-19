package agrpc

import (
	"context"

	"github.com/lightmen/nami/middleware"
	acontext "github.com/lightmen/nami/pkg/acontext"
	"github.com/lightmen/nami/transport"
	"google.golang.org/grpc"
	grpcmd "google.golang.org/grpc/metadata"
)

// unaryServerInterceptor is a gRPC unary server interceptor
func (s *Server) unaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		ctx, cancel := acontext.Merge(ctx, s.baseCtx)
		defer cancel()

		md, _ := grpcmd.FromIncomingContext(ctx)

		header := headerCarrier(md)
		op := GetCmd(ctx, header)
		tr := &Transport{
			operation: op,
			header:    header,
		}
		if s.endpoint != nil {
			tr.endpoint = s.endpoint.String()
		}
		ctx = transport.NewServerContext(ctx, tr)

		if s.timeout > 0 {
			ctx, cancel = context.WithTimeout(ctx, s.timeout)
			defer cancel()
		}

		h := func(ctx context.Context, req any) (any, error) {
			return handler(ctx, req)
		}

		h = middleware.Chain(s.middlewares...)(h)

		reply, err := h(ctx, req)

		return reply, err
	}
}
