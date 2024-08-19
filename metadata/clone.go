package metadata

import (
	"context"
)

func CloneToServerContext(src context.Context) (dst context.Context) {
	dst = context.Background()
	md := New()
	if cliMD, ok := FromClientContext(src); ok {
		md = cliMD.Clone()
	}

	if srvMD, ok := FromServerContext(src); ok {
		//把server的放到后面去Range的原因是： 在server和client的MD有相同key的情况下，
		// 将server的数据 覆盖 client 的数据
		srvMD.Range(func(k, v string) bool {
			md.Set(k, v)
			return true
		})

	}

	return NewServerContext(context.Background(), md)
}

func CloneToClientContext(src context.Context) (dst context.Context) {
	dst = context.Background()
	md := New()
	if srvMD, ok := FromServerContext(src); ok {
		md = srvMD.Clone()
	}

	if cliMD, ok := FromClientContext(src); ok {
		cliMD.Range(func(k, v string) bool {
			md.Set(k, v)
			return true
		})

	}

	return NewClientContext(context.Background(), md)
}
