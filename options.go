package nami

import (
	"context"
	"os"

	"github.com/lightmen/nami/registry"
	"github.com/lightmen/nami/transport"
)

type Option func(*options)

type options struct {
	name      string
	id        string
	servers   []transport.Server
	sigs      []os.Signal
	sigFunc   func(os.Signal) bool
	ctx       context.Context
	metadata  map[string]string
	registrar registry.Registrar
	version   string

	beforeStart []func(context.Context) error
	afterStart  []func(context.Context) error
	afterStop   []func(context.Context) error
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

func Context(ctx context.Context) Option {
	return func(o *options) {
		o.ctx = ctx
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

// Version with service version.
func Version(version string) Option {
	return func(o *options) {
		o.version = version
	}
}

func BeforeStart(fn func(context.Context) error) Option {
	return func(o *options) {
		o.beforeStart = append(o.beforeStart, fn)
	}
}
func AfterStart(fn func(context.Context) error) Option {
	return func(o *options) {
		o.afterStart = append(o.afterStart, fn)
	}
}

func AfterStop(fn func(context.Context) error) Option {
	return func(o *options) {
		o.afterStop = append(o.afterStop, fn)
	}
}

func SigFunc(fn func(os.Signal) bool) Option {
	return func(o *options) {
		o.sigFunc = fn
	}
}
