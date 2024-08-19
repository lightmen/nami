package balancer

import (
	"fmt"

	"github.com/lightmen/nami/pkg/hash/ketama"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
)

type consistentHashPickerBuilder struct {
}

func (b *consistentHashPickerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	if len(info.ReadySCs) == 0 {
		return base.NewErrPicker(balancer.ErrNoSubConnAvailable)
	}

	picker := &consistentHashPicker{
		subConns: make(map[string]balancer.SubConn, len(info.ReadySCs)),
		hash:     ketama.New(),
	}

	for sc, conInfo := range info.ReadySCs {
		node := conInfo.Address.Addr
		picker.hash.Add(node)
		picker.subConns[node] = sc
	}

	return picker
}

type consistentHashPicker struct {
	subConns map[string]balancer.SubConn
	hash     *ketama.Ketama
}

func (p *consistentHashPicker) Pick(info balancer.PickInfo) (result balancer.PickResult, err error) {
	key, ok := FromParamContext(info.Ctx)
	if !ok {
		err = fmt.Errorf("can't found pick value")
		return
	}

	targetAddr, ok := p.hash.Get(key)
	if ok {
		result.SubConn = p.subConns[targetAddr]
	}

	return
}

func (p *consistentHashPicker) Name() string {
	return Consistent
}
