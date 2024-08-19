package agrpc

import (
	"context"
	"net"
	"net/url"
	"time"

	"github.com/lightmen/nami/alog"
	"github.com/lightmen/nami/internal/host"
	"github.com/lightmen/nami/message"
	"github.com/lightmen/nami/middleware"
	"github.com/lightmen/nami/pkg/endpoint"
	"github.com/lightmen/nami/schedule"
	"github.com/lightmen/nami/service"
	"github.com/lightmen/nami/service/cmd"
	"github.com/lightmen/nami/transport"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

var (
	_ transport.Server = (*Server)(nil)
)

type Server struct {
	*grpc.Server
	baseCtx     context.Context
	network     string
	address     string
	lis         net.Listener
	endpoint    *url.URL
	timeout     time.Duration
	middlewares []middleware.Middleware
	unaryInts   []grpc.UnaryServerInterceptor
	streamInts  []grpc.StreamServerInterceptor
	grpcOpts    []grpc.ServerOption
	health      *health.Server

	msgServer message.MessageServer

	sched   schedule.Scheduler
	service service.Service
}

func New(opts ...ServerOption) (srv *Server, err error) {
	srv = &Server{
		network:     "tcp",
		address:     ":0",
		health:      health.NewServer(),
		service:     cmd.GetDefault(),
		timeout:     3 * time.Second,
		middlewares: []middleware.Middleware{},
	}

	for _, opt := range opts {
		opt(srv)
	}

	err = srv.listen()
	if err != nil {
		return srv, err
	}

	unaryInts := []grpc.UnaryServerInterceptor{
		srv.unaryServerInterceptor(),
	}

	streamInts := []grpc.StreamServerInterceptor{}
	if len(srv.unaryInts) > 0 {
		unaryInts = append(unaryInts, srv.unaryInts...)
	}
	if len(srv.streamInts) > 0 {
		streamInts = append(streamInts, srv.streamInts...)
	}
	grpcOpts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(unaryInts...),
		grpc.ChainStreamInterceptor(streamInts...),
	}

	if len(srv.grpcOpts) > 0 {
		grpcOpts = append(grpcOpts, srv.grpcOpts...)
	}
	srv.Server = grpc.NewServer(grpcOpts...)

	if srv.msgServer != nil {
		message.RegisterMessageServer(srv.Server, srv.msgServer)
	} else {
		message.RegisterMessageServer(srv.Server, srv)
	}

	grpc_health_v1.RegisterHealthServer(srv.Server, srv.health)
	reflection.Register(srv.Server)

	return
}

func (s *Server) Start(ctx context.Context) (err error) {
	alog.InfoCtx(ctx, "[gRPC] server lintening on: %s", s.lis.Addr().String())

	s.baseCtx = ctx

	s.health.Resume()

	err = s.Serve(s.lis)

	if err != nil {
		return
	}

	return
}

func (s *Server) Stop(ctx context.Context) (err error) {
	alog.InfoCtx(ctx, "[gRPC] server stopping")

	s.health.Shutdown()
	s.GracefulStop()

	return
}

func (s *Server) listen() error {
	if s.lis != nil {
		return nil
	}

	lis, err := net.Listen(s.network, s.address)
	if err != nil {
		return err
	}

	s.lis = lis

	return nil
}

func (s *Server) Endpoint() (*url.URL, error) {
	if s.endpoint == nil {
		addr, err := host.Extract(s.address, s.lis)
		if err != nil {
			return nil, err
		}

		s.endpoint = endpoint.New(s.Name(), addr)
	}

	return s.endpoint, nil
}

func (s *Server) Name() string {
	return transport.GRPC
}
