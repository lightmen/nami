package service

import (
	"context"

	"github.com/lightmen/nami/codec"
	"github.com/lightmen/nami/middleware"
)

type Handler struct {
	Req    codec.Codec
	Handle middleware.Handler
}

type Registrar interface {
	Register(key any, req codec.Codec, handler middleware.Handler) (err error)
	Get(key any) (handler *Handler, ok bool)
}

type Service interface {
	HandlePacket(ctx context.Context, cmd any, in []byte) (out []byte, err error)
}

var gRegistrar Registrar

func Get() Registrar {
	return gRegistrar
}

func Set(reg Registrar) {
	gRegistrar = reg
}
