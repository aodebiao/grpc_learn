package main

import (
	"context"
	"fmt"
	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"hello_server/pb"
	"log"
	"net"
)

type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, req *pb.HelloRequest) (resp *pb.HelloResponse, err error) {
	reply := fmt.Sprintf("hello,%s", req.Name)
	log.Printf("收到客户端:%s的问候", req.Name)
	return &pb.HelloResponse{Reply: reply}, nil

}

const serverName = "hello_server"

func main() {
	listen, err := net.Listen("tcp", ":7777")
	if err != nil {
		log.Fatalf("err:%v", err)
		return
	}

	grpcServer := grpc.NewServer()
	// 注册服务
	pb.RegisterGreeterServer(grpcServer, &server{})
	// 注册健康检查,增加健康检查的处理逻辑
	// consul发来健康检查请求时,这个负责处理
	healthpb.RegisterHealthServer(grpcServer, health.NewServer())

	// 连接consul
	cc, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		log.Fatalf("err:%v", err)
		return
	}
	// 获取本机出口ip
	ipinfo, err := GetOutboundIP()
	if err != nil {
		log.Fatalf("获取ip失败,err:%v", err)
		return
	}
	log.Printf("获取到的ip:%s", ipinfo.String())
	// 将服务注册到consul
	// 1.定义服务信息
	// 配置健康检查策略,告诉consul如何进行健康检查
	check := &api.AgentServiceCheck{
		GRPC:                           fmt.Sprintf("%s:%d", ipinfo.String(), 7777),
		Timeout:                        "5s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "40s", // 失败后超过40s ,取消注册
	}

	srv := &api.AgentServiceRegistration{
		ID:      fmt.Sprintf("%s-%s-%d", serverName, ipinfo.String(), 7777),
		Name:    serverName,
		Tags:    []string{"aodeibiao"},
		Port:    8888,
		Address: "127.0.0.1",
		Check:   check,
	}
	// 2.注册到consul
	cc.Agent().ServiceRegister(srv)
	// 启动服务
	err = grpcServer.Serve(listen)
	if err != nil {
		log.Fatalf("err:%v", err)
		return
	}

}

func GetOutboundIP() (net.IP, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP, err
}
