package metadata

import "context"

func FromGateContext(ctx context.Context) (string, bool) {
	md, ok := FromServerContext(ctx)
	if !ok {
		return "", false
	}

	addr := md.Get(GateAddrKey)
	if addr == "" {
		return "", false
	}

	return addr, true
}
