package registry

import "context"

//Instance 服务节点信息
type Instance struct {
	ID       string            `json:"ID,omitempty"`
	Name     string            `json:"Name,omitempty"`
	MetaData map[string]string `json:"MetaData,omitempty"`
	// Endpoints is endpoint addresses of the service instance.
	// schema:
	//   http://127.0.0.1:8000?isSecure=false
	//   grpc://127.0.0.1:9000?isSecure=false
	Endpoints []string `json:"Endpoints,omitempty"`
}

type Registrar interface {
	Register(ctx context.Context, service *Instance) error
	Deregister(ctx context.Context, service *Instance) error
}
