package cmd

import (
	"context"
	"fmt"
	"math/rand"
	"reflect"

	"github.com/lightmen/nami/alog"
	"github.com/lightmen/nami/codec"
	"github.com/lightmen/nami/middleware"
	"github.com/lightmen/nami/service"
)

func init() {
	gService = New()

	service.Set(GetDefault())
}

var gService *Service

func GetDefault() *Service {
	return gService
}

type Service struct {
	cmds       map[any]*service.Handler
	middleware []middleware.Middleware
	id         int64
}

func Use(opts ...Option) {
	GetDefault().Use(opts...)
}

func New(opts ...Option) *Service {
	svc := &Service{
		cmds:       make(map[any]*service.Handler),
		id:         rand.Int63(),
		middleware: []middleware.Middleware{
			// tracing.Server(),
			// metadata.Server(),
		},
	}

	for _, o := range opts {
		o(svc)
	}

	return svc
}

func (s *Service) Use(opts ...Option) {
	for _, o := range opts {
		o(s)
	}
}

func (s *Service) Register(cmd any, req codec.Codec, handler middleware.Handler) (err error) {
	_, ok := s.cmds[cmd]
	if ok {
		err = fmt.Errorf("duplicate register cmd: %v", cmd)
		return
	}

	s.cmds[cmd] = &service.Handler{
		Req:    req,
		Handle: handler,
	}

	return
}

func (s *Service) Get(cmd any) (handler *service.Handler, ok bool) {
	proc, ok := s.cmds[cmd]
	if !ok {
		return
	}

	rt := reflect.TypeOf(proc.Req)
	rv := reflect.New(rt.Elem())
	req := rv.Interface().(codec.Codec)

	handler = &service.Handler{
		Req:    req,
		Handle: proc.Handle,
	}

	ok = true
	return
}

func (s *Service) HandlePacket(ctx context.Context, cmd any, in []byte) (out []byte, err error) {
	handler, ok := s.Get(cmd)
	if !ok {
		err = fmt.Errorf("unkown cmd: %d", cmd)
		alog.ErrorCtx(ctx, "%v|s.Get error: %s", cmd, err.Error())
		return
	}

	req := handler.Req
	err = req.Unmarshal(in)
	if err != nil {
		alog.ErrorCtx(ctx, "%v|req.Unmarshal error: %s", cmd, err.Error())
		return
	}

	handle := handler.Handle
	if len(s.middleware) > 0 {
		handle = middleware.Chain(s.middleware...)(handle)
	}

	data, err := handle(ctx, req)

	defer func() { //确保返回前一定会返回
		out = mustDecodeData(ctx, cmd, data)
	}()

	if err != nil {
		alog.ErrorCtx(ctx, "%v|handle error: %s", cmd, err.Error())
		return
	}

	return
}

func mustDecodeData(ctx context.Context, cmd, data any) (out []byte) {
	if data == nil {
		return
	}

	rsp, ok := data.(codec.Codec)
	if !ok {
		// err = aerror.New(codes.InvalidArgument, "rsp is not Codec data")
		alog.ErrorCtx(ctx, "%v|rsp is not Codec data", cmd)
		return
	}

	out, err := rsp.Marshal()
	if err != nil {
		alog.ErrorCtx(ctx, "%v|rsp.Marshal error: %s", cmd, err.Error())
		return
	}

	return
}
