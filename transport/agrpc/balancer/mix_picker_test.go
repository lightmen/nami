package balancer

import (
	"testing"

	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/resolver"
)

func Test_pickerBuildInfoWrapper_String(t *testing.T) {

	subInfo := base.SubConnInfo{
		Address: resolver.Address{
			Addr: "127.0.0.1",
		},
	}
	info := base.PickerBuildInfo{
		ReadySCs: map[balancer.SubConn]base.SubConnInfo{
			nil: subInfo,
		},
	}

	t.Logf("sub: %v", PickerBuildInfoString(info))
}
