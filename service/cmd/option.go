package cmd

import "github.com/lightmen/nami/middleware"

type Option func(s *Service)

func Middleware(m ...middleware.Middleware) Option {
	return func(s *Service) {
		s.middleware = append(s.middleware, m...)
	}
}
