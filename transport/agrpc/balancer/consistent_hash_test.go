package balancer

import (
	"context"
	"testing"

	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/resolver"
)

type mockConn struct{}

func (mc *mockConn) UpdateAddresses([]resolver.Address) {

}

func (mc *mockConn) Connect() {

}

func (mc *mockConn) GetOrBuildProducer(balancer.ProducerBuilder) (p balancer.Producer, close func()) {
	return
}

func Test_consistentHashPicker(t *testing.T) {
	b := &consistentHashPickerBuilder{}

	mc := &mockConn{}

	info := base.PickerBuildInfo{
		ReadySCs: map[balancer.SubConn]base.SubConnInfo{
			mc: {
				Address: resolver.Address{
					Addr: "127.0.0.1:12454",
				},
			},
		},
	}
	picker := b.Build(info)

	pinfo := balancer.PickInfo{
		Ctx: NewParamContext(context.Background(), "testKey1"),
	}

	result, err := picker.Pick(pinfo)
	if err != nil {
		t.Fatalf("Pick got error: %v", err)
	}

	if result.SubConn != mc {
		t.Fatalf("expect %v, got: %v", mc, result.SubConn)
	}

	_, err = picker.Pick(balancer.PickInfo{
		Ctx: context.Background(),
	})

	if err == nil {
		t.Fatalf("expect error, but error is nil")
	}
}
