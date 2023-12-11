package main

import (
	"context"
	"flag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"test/proto"
	"time"
)

// grpc客户端
// 调用server端的 SayHello方法
var name = flag.String("name", "aodebiao", "通过-name告诉server你是谁")

func main() {
	flag.Parse()

	// 连接
	conn, err := grpc.Dial("127.0.0.1:8972", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("grpc.Dial failed,err:%v", err)
	}

	defer conn.Close()
	c := proto.NewGreeterClient(conn) // 使用生成的代码
	// 调用RPC方法
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	resp, err := c.SayHello(ctx, &proto.HelloRequest{Name: *name})
	if err != nil {
		log.Printf("c.SayHello failed,err:%v", err)
		return
	}
	// 拿到结果
	log.Printf("resp value is :%v", resp.Name)
}
