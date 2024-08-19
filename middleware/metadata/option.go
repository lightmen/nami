package metadata

import (
	"strings"

	"github.com/lightmen/nami/metadata"
)

type Option func(o *options)

type options struct {
	prefix []string
	md     metadata.Metadata
}

func (o *options) hasPrefix(key string) bool {
	k := strings.ToLower(key)
	for _, prefix := range o.prefix {
		if strings.HasPrefix(k, prefix) {
			return true
		}
	}
	return false
}

func WithPrefix(prefix string) Option {
	return func(o *options) {
		o.prefix = append(o.prefix, prefix)
	}
}

func WithMetadata(md map[string]string) Option {
	return func(o *options) {
		o.md = md
	}
}
