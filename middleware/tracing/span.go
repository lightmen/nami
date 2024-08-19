package tracing

import (
	"context"
	"net"

	"github.com/lightmen/nami/metadata"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.18.0"
	"go.opentelemetry.io/otel/trace"
)

func setServerSpan(ctx context.Context, span trace.Span) {
	var (
		attrs []attribute.KeyValue
	)

	if md, ok := metadata.FromServerContext(ctx); ok {
		attrs = append(attrs, semconv.PeerServiceKey.String(md.Get(serviceHeader)))
	}

	span.SetAttributes(attrs...)
}

// peerAttr returns attributes about the peer address.
func peerAttr(addr string) []attribute.KeyValue {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return []attribute.KeyValue(nil)
	}

	if host == "" {
		host = "127.0.0.1"
	}

	return []attribute.KeyValue{
		semconv.NetSockPeerAddrKey.String(host),
		semconv.NetSockPeerPortKey.String(port),
	}
}
