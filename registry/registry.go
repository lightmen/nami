package registry

import "context"

//Instance 服务节点信息
type Instance struct {
	ID       string            `json:"ID,omitempty"`
	Name     string            `json:"Name,omitempty"`
	MetaData map[string]string `json:"MetaData,omitempty"`
	// Endpoints is endpoint addresses of the service instance.
	// schema:
	//   http://127.0.0.1:8000
	//   grpc://127.0.0.1:9000
	Endpoints []string `json:"Endpoints,omitempty"`
}

type Registrar interface {
	Register(ctx context.Context, service *Instance) error
	Unregister(ctx context.Context, service *Instance) error
}

// Discovery is service discovery.
type Discovery interface {
	// Watch creates a watcher according to the service name.
	Watch(ctx context.Context, srvName string) error
	// GetService return the service instances in memory according to the service name.
	GetService(ctx context.Context, srvName string) ([]*Instance, error)
}
