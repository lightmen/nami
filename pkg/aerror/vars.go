package aerror

import "github.com/lightmen/nami/codes"

var (
	ServerInternal   = New(codes.Internal, "internal server")           // 内部服务错误
	InvalidParam     = New(codes.InvalidArgument, "invalid param")      //无效参数
	UnsupportMessage = New(codes.Unimplemented, "unsupport message")    //不支持的消息类型
	OutOfRange       = New(codes.OutOfRange, "out of Range")            //超过限制
	DeadlineExceeded = New(codes.DeadlineExceeded, "deadline exceeded") //超时
	NotFound         = New(codes.NotFound, "not found")                 //没有找到
	PermissionDenied = New(codes.PermissionDenied, "permission denied") //没有权限

	Unknown = New(9999, "unknown error") // 不明确的错误
)
