package grpc

import (
	"net"
	"net/url"

	"github.com/lightmen/nami/core/log"
	"github.com/lightmen/nami/internal/endpoint"
	"github.com/lightmen/nami/internal/host"
	"github.com/lightmen/nami/transport"
	"google.golang.org/grpc"
)

type Server struct {
	*grpc.Server
	network  string
	address  string
	lis      net.Listener
	log      log.Logger
	endpoint *url.URL
}

func New(opts ...Option) (srv *Server, err error) {
	srv = &Server{
		network: "tcp",
		address: ":0",
		log:     log.Default(),
	}

	for _, opt := range opts {
		opt(srv)
	}

	err = srv.listen()
	if err != nil {
		return srv, err
	}

	srv.Server = grpc.NewServer() // TODO： need add grpc.ServerOption

	return
}

func (s *Server) Start() (err error) {
	err = s.Serve(s.lis)
	if err != nil {
		return
	}

	s.log.Info("[gRPC] server lintening on: %s", s.lis.Addr().String())
	return
}

func (s *Server) Stop() (err error) {
	s.log.Info("[gRPC] server stopping")
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

		s.endpoint = endpoint.New(transport.GRPC, addr)
	}

	return s.endpoint, nil
}
