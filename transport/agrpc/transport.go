package agrpc

import (
	"context"

	"github.com/lightmen/nami/metadata"
	"github.com/lightmen/nami/transport"
	grpcmd "google.golang.org/grpc/metadata"
)

var _ transport.Transporter = (*Transport)(nil)

// Transport is a gRPC transport.
type Transport struct {
	endpoint  string
	operation string
	header    headerCarrier
}

func GetCmd(ctx context.Context, header transport.Header) string {
	if header == nil {
		var ok bool
		header, ok = metadata.FromClientContext(ctx)
		if !ok {
			return ""
		}
	}

	if header == nil {
		return ""
	}

	return header.Get(metadata.CmdKey)
}

// Kind returns the transport kind.
func (tr *Transport) Kind() transport.Kind {
	return transport.GRPC
}

// Endpoint returns the transport endpoint.
func (tr *Transport) Endpoint() string {
	return tr.endpoint
}

// Operation returns the transport operation.
func (tr *Transport) Operation() string {
	return tr.operation
}

// RequestHeader returns the request header.
func (tr *Transport) Header() transport.Header {
	return tr.header
}

type headerCarrier grpcmd.MD

// Get returns the value associated with the passed key.
func (mc headerCarrier) Get(key string) string {
	vals := grpcmd.MD(mc).Get(key)
	if len(vals) > 0 {
		return vals[0]
	}
	return ""
}

// Set stores the key-value pair.
func (mc headerCarrier) Set(key string, value string) {
	grpcmd.MD(mc).Set(key, value)
}

// Keys lists the keys stored in this carrier.
func (mc headerCarrier) Keys() []string {
	keys := make([]string, 0, len(mc))
	for k := range grpcmd.MD(mc) {
		keys = append(keys, k)
	}
	return keys
}
