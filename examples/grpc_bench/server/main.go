package main

import (
	_ "net/http/pprof"

	"golang.org/x/net/context"

	"github.com/lightmen/nami"
	pb "github.com/lightmen/nami/examples/grpc_bench/helloworld"
	"github.com/lightmen/nami/transport/agrpc"
)

const (
	port = ":50051"
)

// server is used to implement helloworld.GreeterServer.
type server struct{}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	// fmt.Println("######### get client request name :"+in.Name)
	return &pb.HelloReply{Message: "main: Hello " + in.Name}, nil
}

func main() {
	gsrv, err := agrpc.New()
	if err != nil {
		panic(err)
	}

	pb.RegisterGreeterServer(gsrv.Server, &server{})

	app, err := nami.New(
		nami.Name("server"),
		nami.Servers(gsrv),
	)
	if err != nil {
		panic(err)
	}

	if err = app.Run(); err != nil {
		panic(err)
	}
}
