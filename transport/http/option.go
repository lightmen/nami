package http

import (
	"net"
	"net/http"

	"github.com/lightmen/nami/core/chain"
	"github.com/lightmen/nami/core/log"
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
	return func(s *Server) {
		if lis != nil {
			s.lis = lis
		}
	}
}

func Chain(chain chain.Chain) Option {
	return func(s *Server) {
		s.chain = chain
	}
}

func Middlewares(mids ...chain.Middleware) Option {
	return func(s *Server) {
		s.middlewares = append(s.middlewares, mids...)
	}
}

func Handler(handler http.Handler) Option {
	return func(s *Server) {
		s.handler = handler
	}
}
