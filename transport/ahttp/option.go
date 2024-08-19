package ahttp

import (
	"net"
)

type Option func(s *Server)

func Network(network string) Option {
	return func(s *Server) {
		if network != "" {
			s.network = network
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

func Filters(filters ...FilterFunc) Option {
	return func(s *Server) {
		s.filters = append(s.filters, filters...)
	}
}

func PProf(pf bool) Option {
	return func(s *Server) {
		s.usePprof = pf
	}
}
