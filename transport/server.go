package transport

type Server interface {
	Start() error
	Stop() error
}

//Kind define the type of server
type Kind string

func (k Kind) String() string {
	return string(k)
}

const (
	GRPC = "grpc"
	HTTP = "http"
)
