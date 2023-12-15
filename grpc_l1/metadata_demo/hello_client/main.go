package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"hello_client_md/proto"
	"log"
	"time"
)

var name = flag.String("name", "熊二", "-name指定打招呼的对象")

func main() {
	flag.Parse()
	dial, err := grpc.Dial("127.0.0.1:8888", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("grpc dial error:%v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	client := proto.NewGreeterClient(dial)
	// 携带metadata
	md := metadata.Pairs(
		"token", "app-test-aodeibiao",
	)
	ctx = metadata.NewOutgoingContext(ctx, md)
	var header, trailer metadata.MD
	// 接收header 和trailer的metadata
	hello, err := client.SayHello(ctx, &proto.HelloRequest{
		Name: *name,
	}, grpc.Header(&header), grpc.Trailer(&trailer))
	if err != nil {
		log.Fatalf("client.SayHello error:%v", err)
	}
	fmt.Printf("%#v\n", header)
	fmt.Printf("%#v\n", trailer)
	log.Printf("server reply: %v\n", hello.Reply)
	fmt.Printf("%#v\n", header)
	fmt.Printf("%#v\n", trailer)
}
