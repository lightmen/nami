package nami

import (
	"context"
	"os"

	"github.com/lightmen/nami/core/log"
	"github.com/lightmen/nami/registry"
	"github.com/lightmen/nami/transport"
)

type Option func(*options)

type options struct {
	name      string
	id        string
	logger    log.Logger
	servers   []transport.Server
	sigs      []os.Signal
	ctx       context.Context
	metadata  map[string]string
	registrar registry.Registrar
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

func MetaData(metadata map[string]string) Option {
	return func(o *options) {
		o.metadata = metadata
	}
}

func Registrar(registrar registry.Registrar) Option {
	return func(o *options) {
		o.registrar = registrar
	}
}
