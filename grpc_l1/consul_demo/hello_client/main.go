package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/hashicorp/consul/api"
	_ "github.com/mbobakov/grpc-consul-resolver" // 导入
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"hello_client/pb"
	"log"
	"time"
)

var name = flag.String("name", "熊二", "-name指定")

func main() {
	flag.Parse()
	// 1. 连接consul
	/*
		cc, err := api.NewClient(api.DefaultConfig())
		if err != nil {
			log.Fatalf("err:%v", err)
			return
		}
		// 2.根据服务名查询服务实例
		serviceMap, err := cc.Agent().ServicesWithFilter(`Service=="hello_server"`) // 查询服务名称是hello_server的
		if err != nil {
			log.Fatalf("err:%v", err)
			return
		}
		var addr string
		for k, v := range serviceMap {
			fmt.Printf("server key:%s,v:%s", k, v)
			addr = fmt.Sprintf("%s:%d", v.Address, v.Port)
		}
		// 3.与consul返回的服务实例建立连接

		// 4.发起rpc调用
	*/
	conn, err := grpc.Dial("consul://localhost:8500/hello_server?healthy=true",
		//round_robin 作为负载策略
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("err:%v", err)
		return
	}
	defer conn.Close()

	client := pb.NewGreeterClient(conn)

	for i := 0; i < 10; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		resp, err := client.SayHello(ctx, &pb.HelloRequest{Name: *name})
		if err != nil {
			log.Fatalf("err:%v", err)
			return
		}
		log.Printf("收到服务端回复:%s", resp.Reply)

	}

}

// funcName 使用原始方法去consul上查询对应的服务调用再来调用
func funcName() {
	flag.Parse()
	// 1. 连接consul

	cc, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		log.Fatalf("err:%v", err)
		return
	}
	// 2.根据服务名查询服务实例
	serviceMap, err := cc.Agent().ServicesWithFilter(`Service=="hello_server"`) // 查询服务名称是hello_server的
	if err != nil {
		log.Fatalf("err:%v", err)
		return
	}
	var addr string
	for k, v := range serviceMap {
		fmt.Printf("server key:%s,v:%s", k, v)
		addr = fmt.Sprintf("%s:%d", v.Address, v.Port)
	}
	// 3.与consul返回的服务实例建立连接

	// 4.发起rpc调用

	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("err:%v", err)
		return
	}
	defer conn.Close()

	client := pb.NewGreeterClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	resp, err := client.SayHello(ctx, &pb.HelloRequest{Name: *name})
	if err != nil {
		log.Fatalf("err:%v", err)
		return
	}
	log.Printf("收到服务端回复:%s", resp.Reply)
}
