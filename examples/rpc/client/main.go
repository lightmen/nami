package main

import (
	"context"
	"flag"
	"log"
	"math/rand"

	"github.com/lightmen/nami/examples/rpc/internal/api"
	"github.com/lightmen/nami/examples/rpc/internal/rcontext"
	"github.com/lightmen/nami/examples/rpc/internal/reg"
	"github.com/lightmen/nami/pkg/arpc"
	grpccall "github.com/lightmen/nami/pkg/arpc/grpc"
	"github.com/lightmen/nami/pkg/cast"
	"github.com/lightmen/nami/pkg/endpoint"
)

var appName string
var addr string
var count int
var mode string

const (
	Discovery = "discovery"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Llongfile)

	flag.StringVar(&appName, "a", "testSrv", "app name, default is testSrv")
	flag.StringVar(&addr, "d", "127.0.0.1:2379", "etcd addr")
	flag.StringVar(&mode, "m", Discovery, "test mode")
	flag.IntVar(&count, "c", 2, "count")
	flag.Parse()

	_, err := reg.GetRegistry(addr)
	if err != nil {
		panic(err)
	}

	grpccall.Init(arpc.GetDiscorey())

	uidStat := make(map[string]int32)
	nodeStat := make(map[int32]int)

	total := 0
	succ := 0
	startUID := 100000000
	uidList := make([]string, 0)
	for i := 1; i < 1000; i++ {
		uid := cast.ToString(startUID + i)
		uidList = append(uidList, uid)
	}

	target := getTarget(appName, mode)

	log.Printf("got target: %s", target)

	for i := 0; i < count; i++ {
		total++

		idx := rand.Intn(len(uidList))
		uid := uidList[idx]

		req := &api.PingReq{
			UID: uid,
		}

		rsp := &api.PingRsp{}
		actx := rcontext.NewContext(context.Background(), uid)
		err = arpc.Request(actx, target, uid, uid, api.CMD_PING, req, rsp)
		if err != nil {
			log.Printf("arpc.Request for %s error: %s", uid, err.Error())
			panic(err)
		}

		nodeID, ok := uidStat[uid]
		if !ok {
			uidStat[uid] = rsp.ID
		} else if nodeID != rsp.ID {
			log.Printf("%s|warn: node %d not equal %d, rsp: %s", uid, nodeID, rsp.ID, cast.ToJson(rsp))
			break
		}

		nodeStat[rsp.ID]++
		succ++
	}

	log.Printf("stat info, total count: %d, succ: %d, uid count: %d, node stats: %+v\n", total, succ, len(uidStat), nodeStat)
}

func getTarget(appName string, mode string) string {
	if mode == Discovery {
		return appName
	}

	dis := arpc.GetDiscorey()

	instances, err := dis.GetService(context.Background(), appName)
	if err != nil {
		panic(err)
	}

	for _, inst := range instances {
		return endpoint.GetGrpcEndpoint(inst.Endpoints)
	}

	return ""
}
