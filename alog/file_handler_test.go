package alog

import (
	"testing"

	"github.com/lightmen/nami/pkg/cast"
)

func Test_defaultBuildMetadata(t *testing.T) {
	md := defaultBuildMetadata(nil)
	if md == nil {
		md = map[string]any{}
	}
	t.Logf("%s", cast.ToJson(md))
}
