// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.3
// source: cluster.proto

package cluster

import (
	context "context"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// MemberClient is the client API for Member service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MemberClient interface {
	HandleRequest(ctx context.Context, in *RequestMessage, opts ...grpc.CallOption) (*MemberHandleResponse, error)
	HandleNotify(ctx context.Context, in *NotifyMessage, opts ...grpc.CallOption) (*MemberHandleResponse, error)
	HandlePush(ctx context.Context, in *PushMessage, opts ...grpc.CallOption) (*MemberHandleResponse, error)
	HandleResponse(ctx context.Context, in *ResponseMessage, opts ...grpc.CallOption) (*MemberHandleResponse, error)
}

type memberClient struct {
	cc grpc.ClientConnInterface
}

func NewMemberClient(cc grpc.ClientConnInterface) MemberClient {
	return &memberClient{cc}
}

func (c *memberClient) HandleRequest(ctx context.Context, in *RequestMessage, opts ...grpc.CallOption) (*MemberHandleResponse, error) {
	out := new(MemberHandleResponse)
	err := c.cc.Invoke(ctx, "/cluster.Member/HandleRequest", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *memberClient) HandleNotify(ctx context.Context, in *NotifyMessage, opts ...grpc.CallOption) (*MemberHandleResponse, error) {
	out := new(MemberHandleResponse)
	err := c.cc.Invoke(ctx, "/cluster.Member/HandleNotify", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *memberClient) HandlePush(ctx context.Context, in *PushMessage, opts ...grpc.CallOption) (*MemberHandleResponse, error) {
	out := new(MemberHandleResponse)
	err := c.cc.Invoke(ctx, "/cluster.Member/HandlePush", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *memberClient) HandleResponse(ctx context.Context, in *ResponseMessage, opts ...grpc.CallOption) (*MemberHandleResponse, error) {
	out := new(MemberHandleResponse)
	err := c.cc.Invoke(ctx, "/cluster.Member/HandleResponse", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MemberServer is the server API for Member service.
// All implementations should embed UnimplementedMemberServer
// for forward compatibility
type MemberServer interface {
	HandleRequest(context.Context, *RequestMessage) (*MemberHandleResponse, error)
	HandleNotify(context.Context, *NotifyMessage) (*MemberHandleResponse, error)
	HandlePush(context.Context, *PushMessage) (*MemberHandleResponse, error)
	HandleResponse(context.Context, *ResponseMessage) (*MemberHandleResponse, error)
}

// UnimplementedMemberServer should be embedded to have forward compatible implementations.
type UnimplementedMemberServer struct {
}

func (UnimplementedMemberServer) HandleRequest(context.Context, *RequestMessage) (*MemberHandleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HandleRequest not implemented")
}
func (UnimplementedMemberServer) HandleNotify(context.Context, *NotifyMessage) (*MemberHandleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HandleNotify not implemented")
}
func (UnimplementedMemberServer) HandlePush(context.Context, *PushMessage) (*MemberHandleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HandlePush not implemented")
}
func (UnimplementedMemberServer) HandleResponse(context.Context, *ResponseMessage) (*MemberHandleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HandleResponse not implemented")
}

// UnsafeMemberServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MemberServer will
// result in compilation errors.
type UnsafeMemberServer interface {
	mustEmbedUnimplementedMemberServer()
}

func RegisterMemberServer(s grpc.ServiceRegistrar, srv MemberServer) {
	s.RegisterService(&Member_ServiceDesc, srv)
}

func _Member_HandleRequest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MemberServer).HandleRequest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cluster.Member/HandleRequest",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MemberServer).HandleRequest(ctx, req.(*RequestMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _Member_HandleNotify_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NotifyMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MemberServer).HandleNotify(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cluster.Member/HandleNotify",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MemberServer).HandleNotify(ctx, req.(*NotifyMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _Member_HandlePush_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PushMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MemberServer).HandlePush(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cluster.Member/HandlePush",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MemberServer).HandlePush(ctx, req.(*PushMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _Member_HandleResponse_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ResponseMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MemberServer).HandleResponse(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cluster.Member/HandleResponse",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MemberServer).HandleResponse(ctx, req.(*ResponseMessage))
	}
	return interceptor(ctx, in, info, handler)
}

// Member_ServiceDesc is the grpc.ServiceDesc for Member service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Member_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "cluster.Member",
	HandlerType: (*MemberServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "HandleRequest",
			Handler:    _Member_HandleRequest_Handler,
		},
		{
			MethodName: "HandleNotify",
			Handler:    _Member_HandleNotify_Handler,
		},
		{
			MethodName: "HandlePush",
			Handler:    _Member_HandlePush_Handler,
		},
		{
			MethodName: "HandleResponse",
			Handler:    _Member_HandleResponse_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "cluster.proto",
}
