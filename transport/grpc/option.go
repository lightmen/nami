package grpc

import (
	"net"

	"github.com/lightmen/nami/log"
)

type Option func(s *Server)

func Address(addr string) Option {
	return func(s *Server) {
		if addr != "" {
			s.address = addr
		}
	}
}

func Network(network string) Option {
	return func(s *Server) {
		if network != "" {
			s.network = network
		}
	}
}

func Log(log log.Logger) Option {
	return func(s *Server) {
		if log != nil {
			s.log = log
		}
	}
}

func Listen(lis net.Listener) Option {

}
