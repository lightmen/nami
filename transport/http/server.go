package http

import (
	"context"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/lightmen/nami/core/chain"
	"github.com/lightmen/nami/core/log"
	"github.com/lightmen/nami/internal/endpoint"
	"github.com/lightmen/nami/internal/host"
	"github.com/lightmen/nami/transport"
	"github.com/lightmen/nami/transport/http/handler"
)

type Server struct {
	*http.Server
	lis         net.Listener
	network     string
	address     string
	timeout     time.Duration
	log         log.Logger
	middlewares []chain.Middleware
	chain       chain.Chain
	handler     http.Handler
	endpoint    *url.URL
}

func New(opts ...Option) (s *Server, err error) {
	s = &Server{
		network:     "tcp",
		address:     ":0",
		timeout:     2 * time.Second,
		log:         log.Default(),
		middlewares: []chain.Middleware{},
		chain:       chain.New(handler.Recover),
	}

	for _, opt := range opts {
		opt(s)
	}

	if len(s.middlewares) > 0 {
		s.chain.Append(s.middlewares...)
	}

	s.Server = &http.Server{
		Handler: s.chain.Then(s.handler),
	}

	err = s.listen()
	if err != nil {
		return nil, err
	}

	return
}

func (s *Server) Start(ctx context.Context) (err error) {
	err = s.Serve(s.lis)
	if err != nil {
		return
	}

	s.log.Info("[HTTP] server lintening on: %s", s.lis.Addr().String())
	return
}

func (s *Server) Stop(ctx context.Context) (err error) {
	s.log.Info("[HTTP] server stopping")
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
		if err := s.listen(); err != nil {
			return nil, err
		}

		addr, err := host.Extract(s.address, s.lis)
		if err != nil {
			return nil, err
		}

		s.endpoint = endpoint.New(transport.HTTP, addr)
	}

	return s.endpoint, nil
}
