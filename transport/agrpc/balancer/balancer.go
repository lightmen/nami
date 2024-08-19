package balancer

import (
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
)

func init() {
	//默认使用一致性hash
	b := base.NewBalancerBuilder(Mix, &mixPickerBuilder{}, base.Config{HealthCheck: true})
	balancer.Register(b)
}
