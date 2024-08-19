package nami

import (
	"context"

	"github.com/lightmen/nami/transport"
)

type AppInfo interface {
	ID() string
	Name() string
	Version() string
	Metadata() map[string]string
	Endpoint() []string
	Server(name string) transport.Server
}

type appKey struct{}

// NewContext returns a new Context that carries value.
func NewContext(ctx context.Context, s AppInfo) context.Context {
	return context.WithValue(ctx, appKey{}, s)
}

// FromContext returns the Transport value stored in ctx, if any.
func FromContext(ctx context.Context) (s AppInfo, ok bool) {
	s, ok = ctx.Value(appKey{}).(AppInfo)
	return
}

func (a *App) Name() string {
	return a.opts.name
}

func (a *App) ID() string {
	return a.opts.id
}

// Version returns app version.
func (a *App) Version() string { return a.opts.version }

// Metadata returns service metadata.
func (a *App) Metadata() map[string]string { return a.opts.metadata }

func (a *App) Endpoint() []string {
	if a.instance != nil {
		return a.instance.Endpoints
	}
	return nil
}

// Server 返回
func (a *App) Server(name string) transport.Server {
	for _, srv := range a.opts.servers {
		if srv.Name() == name {
			return srv
		}
	}

	return nil
}
