package balancer

import "context"

type nameKey struct{}
type targetKey struct{}
type paramKey struct{}

// NewNameContext 设置balance负载的名字, mix负载会根据该名字选择其他balance
func NewNameContext(ctx context.Context, name string) context.Context {
	ctx = context.WithValue(ctx, nameKey{}, name)
	return ctx
}

func FromNameContext(ctx context.Context) (string, bool) {
	name, ok := ctx.Value(nameKey{}).(string)
	return name, ok
}

// NewParamContext 用于balance负载的参数
func NewParamContext(ctx context.Context, addr string) context.Context {
	return context.WithValue(ctx, paramKey{}, addr)
}

func FromParamContext(ctx context.Context) (string, bool) {
	addr, ok := ctx.Value(paramKey{}).(string)
	return addr, ok
}

// NewTargetContext 设置target目标服务的名字
func NewTargetContext(ctx context.Context, name string) context.Context {
	ctx = context.WithValue(ctx, targetKey{}, name)
	return ctx
}

func FromTargetContext(ctx context.Context) (string, bool) {
	name, ok := ctx.Value(targetKey{}).(string)
	return name, ok
}
