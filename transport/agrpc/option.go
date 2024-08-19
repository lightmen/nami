package agrpc

import (
	"net"
	"time"

	"github.com/lightmen/nami/message"
	"github.com/lightmen/nami/middleware"
	"github.com/lightmen/nami/schedule"
	"github.com/lightmen/nami/service"
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

func Listen(lis net.Listener) ServerOption {
	return func(s *Server) {
		if lis != nil {
			s.lis = lis
		}
	}
}

func MessageServer(ms message.MessageServer) ServerOption {
	return func(s *Server) {
		s.msgServer = ms
	}
}

func Scheduler(sched schedule.Scheduler) ServerOption {
	return func(s *Server) {
		s.sched = sched
	}
}

func Service(svc service.Service) ServerOption {
	return func(s *Server) {
		s.service = svc
	}
}

func Timeout(timeout time.Duration) ServerOption {
	return func(s *Server) {
		s.timeout = timeout
	}
}

func Middleware(m ...middleware.Middleware) ServerOption {
	return func(s *Server) {
		s.middlewares = append(s.middlewares, m...)
	}
}
