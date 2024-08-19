package arpc

import (
	"context"

	"github.com/lightmen/nami/alog"
	"github.com/lightmen/nami/pkg/endpoint"
	"github.com/lightmen/nami/registry"
)

var gDiscovery registry.Discovery

func SetDiscorey(dis registry.Discovery) {
	gDiscovery = dis
}

func GetDiscorey() registry.Discovery {
	return gDiscovery
}

func MustGetGrpcAddrsByName(ctx context.Context, srvName string) []string {
	dis := GetDiscorey()
	if dis == nil {
		return nil
	}

	insList, err := dis.GetService(ctx, srvName)
	if err != nil {
		alog.ErrorCtx(ctx, "got %s addr error: %s", srvName, err.Error())
		return nil
	}

	addrList := make([]string, 0, len(insList))
	for _, ins := range insList {
		addr := endpoint.GetGrpcEndpoint(ins.Endpoints)
		if addr == "" {
			continue
		}

		addrList = append(addrList, addr)
	}

	return addrList
}
