package main

import (
	"flag"
	"log"
	"math/rand"

	"github.com/lightmen/nami"
	"github.com/lightmen/nami/examples/rpc/internal/api"
	"github.com/lightmen/nami/examples/rpc/internal/reg"
	"github.com/lightmen/nami/transport/agrpc"
)

var appName string
var srvID int

func main() {
	flag.StringVar(&appName, "a", "testSrv", "app name, default is testSrv")
	flag.IntVar(&srvID, "i", 0, "srvID")
	flag.Parse()

	if srvID == 0 {
		srvID = rand.Intn(100)
	}

	log.Printf("start server for %s:%d", appName, srvID)
	reg, err := reg.GetRegistry("")
	if err != nil {
		panic(err)
	}

	gsrv, err := agrpc.New()
	if err != nil {
		panic(err)
	}

	api.Register(int32(srvID))

	app, err := nami.New(
		nami.Name(appName),
		nami.Registrar(reg),
		nami.Servers(gsrv),
	)
	if err != nil {
		panic(err)
	}

	if err = app.Run(); err != nil {
		panic(err)
	}
}
