package balancer

import (
	"fmt"

	"github.com/lightmen/nami/codes"
	"github.com/lightmen/nami/pkg/aerror"
	"github.com/lightmen/nami/pkg/cast"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
)

const (
	Mix        = "mix"        // consistent 和 其他注册balancer算法 的混合体
	Consistent = "consistent" //一致性hash
)

var gSelector *selector

func init() {
	gSelector = &selector{
		pickers: make(map[string]Picker),
	}
}

func RegisterPicker(builder Picker) {
	gSelector.Register(builder)
}

type PickerBuildInfoString base.PickerBuildInfo

func (p PickerBuildInfoString) String() string {
	infos := make([]base.SubConnInfo, 0, len(p.ReadySCs))
	for _, info := range p.ReadySCs {
		infos = append(infos, info)
	}

	return cast.ToJson(infos)
}

type mixPickerBuilder struct {
}

func (b *mixPickerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	if len(info.ReadySCs) == 0 {
		return base.NewErrPicker(balancer.ErrNoSubConnAvailable)
	}

	// alog.Debug("mix picker build got(%d): %s", len(info.ReadySCs), PickerBuildInfoString(info).String())

	builders := map[string]base.PickerBuilder{
		Consistent: &consistentHashPickerBuilder{},
	}

	mixPicker := &mixPicker{
		pickers: make(map[string]balancer.Picker, len(builders)),
	}

	gSelector.Update(info)

	for name, builder := range builders {
		picker := builder.Build(info)
		mixPicker.pickers[name] = picker
	}

	return mixPicker
}

type mixPicker struct {
	pickers map[string]balancer.Picker
}

func (p *mixPicker) Pick(info balancer.PickInfo) (result balancer.PickResult, err error) {
	//1. 首先权重最高的Picker是直连，如果是直连负载，直接使用直连
	//2. 如果业务注册了picker负载，判断是否匹配业务的picker,如果匹配，使用业务的picker
	//3. 否则的话，只用Consistent一致性hash的负载

	ctx := info.Ctx

	//获取默认的picker
	name, ok := FromNameContext(ctx)
	if !ok {
		name = Consistent
	}

	if picker := gSelector.Get(ctx); picker != nil {
		return picker.Pick(info)
	}

	picker, ok := p.pickers[name]
	if !ok {
		err = aerror.New(codes.InvalidArgument, fmt.Sprintf("got empty picker for:%s", name))
		return
	}

	return picker.Pick(info)
}
