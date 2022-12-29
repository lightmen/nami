package nami

import (
	"os"

	"github.com/lightmen/nami/log"
	"github.com/lightmen/nami/transport"
)

type Option func(*options)

type options struct {
	name    string
	id      string
	logger  log.Logger
	servers []transport.Server
	sigs    []os.Signal
}

func Name(name string) Option {
	return func(o *options) {
		o.name = name
	}
}

func ID(id string) Option {
	return func(o *options) {
		if id != "" {
			o.id = id
		}
	}
}

func Logger(logger log.Logger) Option {
	return func(o *options) {
		if logger != nil {
			o.logger = logger
		}
	}
}

func Servers(srv ...transport.Server) Option {
	return func(o *options) {
		o.servers = srv
	}
}
