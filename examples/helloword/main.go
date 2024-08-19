package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/lightmen/nami"
	"github.com/lightmen/nami/transport/ahttp"
)

type count struct {
	val int32
}

func (c *count) Count(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.val++

		next.ServeHTTP(w, r)

		log.Printf("total req count: %d", c.val)
	})
}

func timeStat(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		next.ServeHTTP(w, r)

		takes := time.Since(startTime)
		log.Printf("req takes: %v", takes)
	})
}

var (
	c = &count{}
)

func main() {
	hsrv, err := ahttp.New(
		":9901",
		ahttp.Middlewares(c.Count, timeStat),
	)
	if err != nil {
		panic(err)
	}

	hsrv.Handle("/", &helloHandler{})

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

	rd := rand.Int63n(500)
	time.Sleep(time.Millisecond * time.Duration(rd))
}
