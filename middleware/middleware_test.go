package middleware

import (
	"context"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func tagMiddleware(tag string) Middleware {
	return func(next Handler) Handler {
		return func(ctx context.Context, req any) (any, error) {
			log.Printf("%v", tag)

			rsp, err := next(ctx, req)
			val := tag + rsp.(string)

			return val, err
		}
	}
}

func TestChain(t *testing.T) {
	testHandler := func(ctx context.Context, req any) (rsp any, err error) {
		log.Printf("testHandler: %v", req)
		return req, nil
	}

	t1 := tagMiddleware("t1")
	t2 := tagMiddleware("t2")
	rsp, _ := Chain(t1, t2)(testHandler)(nil, "app")
	assert.Equal(t, "t1t2app", rsp.(string),
		"Then does not order handlers correctly")
}
