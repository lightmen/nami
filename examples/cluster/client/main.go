package main

import (
	"context"
	"log"

	"github.com/lightmen/nami/message"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("127.0.0.1:9902",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	cli := message.NewMessageClient(conn)

	ctx := context.Background()
	packet := &message.Packet{
		Head: &message.Head{
			Route: "100001",
			Type:  message.REQUEST,
		},
	}
	_, err = cli.HandleMessage(ctx, packet)
	if err != nil {
		log.Printf("HandleRequest error: %s", err.Error())
	}
}
