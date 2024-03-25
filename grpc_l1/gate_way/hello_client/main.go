package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
	"way_client1/pb"
)

var name = flag.String("name", "熊二", "-name指定名称")

func main() {
	flag.Parse()
	dial, err := grpc.Dial("127.0.0.1:4567", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("dial error:%v\n", err)
	}
	client := pb.NewGreeterClient(dial)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	resp, err := client.SayHello(ctx, &pb.HelloRequest{
		Name: *name,
	})
	if err != nil {
		log.Fatalf("recv SayHello error:%v\n", err)
	}
	fmt.Printf("recv value:%v\n", resp.Reply)

}
