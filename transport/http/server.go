package http

import (
	"net"
	"net/http"
	"time"

	"github.com/lightmen/nami/log"
)

type Server struct {
	*http.Server
	lis     net.Listener
	network string
	address string
	timeout time.Duration
	log     log.Logger
}

func New(opts ...Option) (s *Server, err error) {
	s = &Server{
		network: "tcp",
		address: ":0",
		timeout: 2 * time.Second,
		log:     log.Default(),
	}

	for _, opt := range opts {
		opt(s)
	}

	s.Server = &http.Server{}

	err = s.listen()
	if err != nil {
		return nil, err
	}

	return
}

func (s *Server) Start() (err error) {
	err = s.Serve(s.lis)
	if err != nil {
		return
	}

	s.log.Info("[HTTP] server lintening on: %s", s.lis.Addr().String())
	return
}

func (s *Server) Stop() (err error) {
	s.log.Info("[HTTP] server stopping")
	return
}

func (s *Server) listen() error {
	if s.lis == nil {
		lis, err := net.Listen(s.network, s.address)
		if err != nil {
			return err
		}

		s.lis = lis
	}
	return nil
}
