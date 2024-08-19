package pool

import (
	"context"
	"time"

	"github.com/lightmen/nami/transport/agrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

const (
	// DialTimeout the timeout of create connection
	DialTimeout = 5 * time.Second

	// BackoffMaxDelay provided maximum delay when backing off after failed connection attempts.
	BackoffMaxDelay = 3 * time.Second

	// KeepAliveTime is the duration of time after which if the client doesn't see
	// any activity it pings the server to see if the transport is still alive.
	KeepAliveTime = time.Duration(10) * time.Second

	// KeepAliveTimeout is the duration of time for which the client waits after having
	// pinged for keepalive check and if no activity is seen even after that the connection
	// is closed.
	KeepAliveTimeout = time.Duration(3) * time.Second

	// InitialWindowSize we set it 1GB is to provide system's throughput.
	InitialWindowSize = 1 << 30

	// InitialConnWindowSize we set it 1GB is to provide system's throughput.
	InitialConnWindowSize = 1 << 30

	// MaxSendMsgSize set max gRPC request message size sent to server.
	// If any request message size is larger than current value, an error will be reported from gRPC.
	MaxSendMsgSize = 4 << 30

	// MaxRecvMsgSize set max gRPC receive message size received from server.
	// If any message size is larger than current value, an error will be reported from gRPC.
	MaxRecvMsgSize = 4 << 30
)

type Option func(o *options)

type options struct {
	// Dial is an application supplied function for creating and configuring a connection.
	Dial func(address string) (*grpc.ClientConn, error)

	// Maximum number of idle connections in the pool.
	MaxIdle int

	// Maximum number of connections allocated by the pool at a given time.
	// When zero, there is no limit on the number of connections in the pool.
	MaxActive int

	// MaxConcurrentStreams limit on the number of concurrent streams to each single connection
	MaxConcurrentStreams int

	// If Reuse is true and the pool is at the MaxActive limit, then Get() reuse
	// the connection to return, If Reuse is false and the pool is at the MaxActive limit,
	// create a one-time connection to return.
	Reuse bool
}

func WithDial(dial func(address string) (*grpc.ClientConn, error)) Option {
	return func(o *options) {
		o.Dial = dial
	}
}

// DefaultOptions sets a list of recommended options for good performance.
// Feel free to modify these to suit your needs.
var DefaultOptions = options{
	Dial:                 Dial,
	MaxIdle:              4,
	MaxActive:            16,
	MaxConcurrentStreams: 16,
	Reuse:                true,
}

// Dial return a grpc connection with defined configurations.
func Dial(address string) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DialTimeout)
	defer cancel()

	return agrpc.Dial(
		ctx,
		agrpc.WithEndpoint(address),
		agrpc.WithOptions(
			grpc.WithBackoffMaxDelay(BackoffMaxDelay),
			grpc.WithInitialWindowSize(InitialWindowSize),
			grpc.WithInitialConnWindowSize(InitialConnWindowSize),
			grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(MaxSendMsgSize)),
			grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(MaxRecvMsgSize)),
			grpc.WithKeepaliveParams(keepalive.ClientParameters{
				Time:                KeepAliveTime,
				Timeout:             KeepAliveTimeout,
				PermitWithoutStream: true,
			})),
	)
}
