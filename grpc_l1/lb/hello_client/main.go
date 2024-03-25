package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"lb_hello_client/pb"
	"log"
	"time"
)

var name = flag.String("name", "王大勇", "-name指定名称")
var port = flag.String("p", "8991", "")

func main() {
	flag.Parse()
	// 自定义解析
	//conn, err := grpc.Dial("aodeibiao:///resolver.aodeibiao.com",
	//	grpc.WithTransportCredentials(insecure.NewCredentials()),
	//	//grpc.WithResolvers(&aodebiaoResolverBuilder{}), // 指定使用q1miResolverBuilder
	//)
	// 自定义解析 + 轮询的负载均衡
	conn, err := grpc.Dial("aodeibiao:///resolver.aodeibiao.com",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
		//grpc.WithResolvers(&aodebiaoResolverBuilder{}), // 指定使用q1miResolverBuilder
	)

	if err != nil {
		log.Fatalf("dial 127.0.0.1:8991 error:%v\n", err)
		return
	}
	defer conn.Close()
	client := pb.NewGreeterClient(conn)
	for i := 0; i < 10; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		hello, err := client.SayHello(ctx, &pb.HelloRequest{
			Name: *name,
		})
		if err != nil {
			fmt.Printf("SayHello error:%v\n", err)
			return
		}
		fmt.Printf("接收到的回应:%s\n", hello.Reply)
	}
}

func dnsResolver() {
	flag.Parse()
	conn, err := grpc.Dial(fmt.Sprintf("dns:///127.0.0.1:%s", *port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("dial 127.0.0.1:8991 error:%v\n", err)
		return
	}
	defer conn.Close()
	client := pb.NewGreeterClient(conn)
	for i := 0; i < 10; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		hello, err := client.SayHello(ctx, &pb.HelloRequest{
			Name: *name,
		})
		if err != nil {
			fmt.Printf("SayHello error:%v\n", err)
			return
		}
		fmt.Printf("接收到的回应:%s\n", hello.Reply)
	}
}
