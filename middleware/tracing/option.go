package tracing

import (
	"strings"

	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type Option func(o *options)

type options struct {
	name       string
	prefix     []string
	propagator propagation.TextMapPropagator
	provider   trace.TracerProvider
}

func (o *options) hasPrefix(key string) bool {
	k := strings.ToLower(key)
	for _, prefix := range o.prefix {
		if strings.HasPrefix(k, prefix) {
			return true
		}
	}
	return false
}

func WithName(name string) Option {
	return func(o *options) {
		o.name = name
	}
}

func WithProvider(provider trace.TracerProvider) Option {
	return func(o *options) {
		o.provider = provider
	}
}

func WithPropagator(propagator propagation.TextMapPropagator) Option {
	return func(o *options) {
		o.propagator = propagator
	}
}
