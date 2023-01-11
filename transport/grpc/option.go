package grpc

import (
	"net"

	"github.com/lightmen/nami/core/log"
)

type ServerOption func(s *Server)

func Address(addr string) ServerOption {
	return func(s *Server) {
		if addr != "" {
			s.address = addr
		}
	}
}

func Network(network string) ServerOption {
	return func(s *Server) {
		if network != "" {
			s.network = network
		}
	}
}

func Log(log log.Logger) ServerOption {
	return func(s *Server) {
		if log != nil {
			s.log = log
		}
	}
}

func Listen(lis net.Listener) ServerOption {
	return func(s *Server) {
		if lis != nil {
			s.lis = lis
		}
	}
}
