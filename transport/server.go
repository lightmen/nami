package transport

import "context"

type Server interface {
	Start(context.Context) error
	Stop(context.Context) error
	Name() string
}

// Kind define the type of server
type Kind string

func (k Kind) String() string {
	return string(k)
}

const (
	GRPC = "grpc"
	HTTP = "http"
)
