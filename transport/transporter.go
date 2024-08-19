package transport

import (
	"context"

	"github.com/lightmen/nami/pkg/acontext"
)

type Header interface {
	Get(key string) string
	Set(key, value string)
	Keys() []string
}

type Transporter interface {
	// Endpoint return server or client endpoint
	// Server Transport: grpc://127.0.0.1:9000
	// Client Transport: discovery:///provider-demo
	Endpoint() string

	// Operation Service full method
	// http: /api/check_update
	// grpc: CMD_GAME_LOGIN
	Operation() string

	Header() Header
}

type (
	serverTransportKey struct{}
	clientTransportKey struct{}
)

func NewServerContext(ctx context.Context, tr Transporter) context.Context {
	return acontext.WithValue(ctx, serverTransportKey{}, tr)
}

func FromServerContext(ctx context.Context) (tr Transporter, ok bool) {
	tr, ok = ctx.Value(serverTransportKey{}).(Transporter)
	return
}

func NewClientContext(ctx context.Context, tr Transporter) context.Context {
	return acontext.WithValue(ctx, clientTransportKey{}, tr)
}

func FromClientContext(ctx context.Context) (tr Transporter, ok bool) {
	tr, ok = ctx.Value(clientTransportKey{}).(Transporter)
	return
}
