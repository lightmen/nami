package session

import (
	"context"

	"github.com/lightmen/nami/pkg/acontext"
)

type sessionKey struct{}

// NewContext 将Session存储到context.Context中
func NewContext(ctx context.Context, s Session) context.Context {
	ctx = acontext.WithValue(ctx, sessionKey{}, s)
	return ctx
}

// FromContext 从context.Context中获取session
func FromContext(ctx context.Context) (Session, bool) {
	s, ok := ctx.Value(sessionKey{}).(Session)
	return s, ok
}
