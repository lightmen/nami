package registry

import "context"

// Instance 服务节点信息
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
	// GetService return the service instances in memory according to the service name. if srvName is emtpy, get all service
	GetService(ctx context.Context, srvName string) ([]*Instance, error)

	// Watch creates a watcher according to the service name, if srvName is emtpy, watch all service
	Watch(ctx context.Context, srvName string) (Watcher, error)
}

// Watcher is service watcher.
type Watcher interface {
	// Next returns services in the following two cases:
	// 1.the first time to watch and the service instance list is not empty.
	// 2.any service instance changes found.
	// if the above two conditions are not met, it will block until context deadline exceeded or canceled
	Next() ([]*Instance, error)
	// Stop close the watcher.
	Stop() error
}
