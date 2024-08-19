package grpc

import "github.com/lightmen/nami/registry"

type clientOption func(*client)

func WithDiscovery(dis registry.Discovery) clientOption {
	return func(c *client) {
		c.dis = dis
	}
}
