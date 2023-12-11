package main

import (
	"aodebiao/pb"
	"context"
	"flag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

var (
	x = flag.Int("x", 1, "通过-x告诉我第一位数")
	y = flag.Int("y", 1, "通过-y告诉我第一位数")
)

func main() {
	flag.Parse()

	conn, err := grpc.Dial("127.0.0.1:8888", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("grpc.Dial failed,err:%v", err)
	}
	defer conn.Close()

	client := pb.NewAddServerClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	response, err := client.Add(ctx, &pb.AddRequest{X: (int32)(*x), Y: (int32)(*y)})
	if err != nil {
		log.Fatalf("client.Add failed,error:%v", err)
	}
	log.Printf("client.Add get result vaule: %v", response.GetResult())
}
