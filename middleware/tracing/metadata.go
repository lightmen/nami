package tracing

import (
	"context"

	"github.com/lightmen/nami"
	"github.com/lightmen/nami/metadata"
	"go.opentelemetry.io/otel/propagation"
)

const serviceHeader = "service-name"

type Metadata struct{}

var _ propagation.TextMapPropagator = Metadata{}

// Inject set cross-cutting concerns from the Context into the carrier.
func (b Metadata) Inject(ctx context.Context, carrier propagation.TextMapCarrier) {
	app, ok := nami.FromContext(ctx)
	if !ok {
		return
	}

	carrier.Set(serviceHeader, app.Name())
}

// Extract reads cross-cutting concerns from the carrier into a Context.
func (b Metadata) Extract(parent context.Context, carrier propagation.TextMapCarrier) context.Context {
	name := carrier.Get(serviceHeader)
	if name == "" {
		return parent
	}

	if md, ok := metadata.FromServerContext(parent); ok {
		md.Set(serviceHeader, name)
		return parent
	}

	md := metadata.New()
	md.Set(serviceHeader, name)
	parent = metadata.NewServerContext(parent, md)
	return parent
}

// Fields returns the keys who's values are set with Inject.
func (b Metadata) Fields() []string {
	return []string{serviceHeader}
}
