package grpc

import (
	"context"

	"github.com/lightmen/nami/internal/cluster"
)

func (s *Server) HandleRequest(ctx context.Context, req *cluster.RequestMessage) (rsp *cluster.MemberHandleResponse, err error) {
	return
}

func (s *Server) HandleNotify(ctx context.Context, req *cluster.NotifyMessage) (rsp *cluster.MemberHandleResponse, err error) {
	return
}

func (s *Server) HandlePush(ctx context.Context, req *cluster.PushMessage) (rsp *cluster.MemberHandleResponse, err error) {
	return
}

func (s *Server) HandleResponse(ctx context.Context, req *cluster.ResponseMessage) (rsp *cluster.MemberHandleResponse, err error) {
	return
}
