package main

import (
	"github.com/lightmen/nami"
	"github.com/lightmen/nami/transport/agrpc"
)

func main() {
	gsrv, err := agrpc.New(
		agrpc.Address("127.0.0.1:9902"),
	)
	if err != nil {
		panic(err)
	}

	app, err := nami.New(
		nami.Servers(gsrv),
	)
	if err != nil {
		panic(err)
	}

	if err = app.Run(); err != nil {
		panic(err)
	}
}
