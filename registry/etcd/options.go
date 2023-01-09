package etcd

import (
	"context"
	"time"

	"github.com/lightmen/nami/core/log"
)

type Option func(o *options)

type options struct {
	ctx       context.Context
	namespace string
	ttl       time.Duration
	log       log.Logger
}

func Context(ctx context.Context) Option {
	return func(o *options) {
		o.ctx = ctx
	}
}

func Namespace(namespace string) Option {
	return func(o *options) {
		o.namespace = namespace
	}
}

func TTL(ttl time.Duration) Option {
	return func(o *options) {
		o.ttl = ttl
	}
}

func Log(log log.Logger) Option {
	return func(o *options) {
		o.log = log
	}
}
