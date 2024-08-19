package metadata

import "context"

func FromCmdContext(ctx context.Context) string {
	md, ok := FromServerContext(ctx)
	if !ok {
		return ""
	}

	addr := md.Get(CmdKey)
	if addr == "" {
		return ""
	}

	return addr
}
