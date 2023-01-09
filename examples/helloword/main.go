package main

import (
	"net/http"

	"github.com/lightmen/nami"
	transhttp "github.com/lightmen/nami/transport/http"
)

func main() {
	hsrv, err := transhttp.New(
		transhttp.Address(":9901"),
		transhttp.Handler(&helloHandler{}),
	)
	if err != nil {
		panic(err)
	}

	app, err := nami.New(
		nami.Servers(hsrv),
	)
	if err != nil {
		panic(err)
	}

	if err = app.Run(); err != nil {
		panic(err)
	}
}

type helloHandler struct{}

func (t *helloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello world"))
}
