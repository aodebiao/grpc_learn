package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"hello_server/pb"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

type server struct {
	pb.UnimplementedGreeterServer
	addr string
}

func (s *server) SayHello(ctx context.Context, req *pb.HelloRequest) (resp *pb.HelloResponse, err error) {
	reply := fmt.Sprintf("%s hello,%s", s.addr, req.Name)
	log.Printf("收到客户端:%s的问候", req.Name)
	return &pb.HelloResponse{Reply: reply}, nil

}

const serverName = "hello_server"

var port = flag.String("port", "8888", "-port")

func main() {
	flag.Parse()
	addr := fmt.Sprintf(":%s", *port)
	fmt.Println("addr:", addr)
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("err:%v", err)
		return
	}

	grpcServer := grpc.NewServer()
	// 注册服务
	// 获取本机出口ip
	ipinfo, err := GetOutboundIP()
	if err != nil {
		log.Fatalf("获取ip失败,err:%v", err)
		return
	}
	pb.RegisterGreeterServer(grpcServer, &server{addr: fmt.Sprintf("%s:%s", ipinfo.String(), *port)})
	// 注册健康检查,增加健康检查的处理逻辑
	// consul发来健康检查请求时,这个负责处理
	healthpb.RegisterHealthServer(grpcServer, health.NewServer())

	// 连接consul
	cc, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		log.Fatalf("err:%v", err)
		return
	}

	log.Printf("获取到的ip:%s", ipinfo.String())
	// 将服务注册到consul
	// 1.定义服务信息
	// 配置健康检查策略,告诉consul如何进行健康检查
	checkAddr := fmt.Sprintf("%s:%s", ipinfo.String(), *port)
	fmt.Println("checkAddr:", checkAddr)
	check := &api.AgentServiceCheck{
		GRPC:                           checkAddr,
		Timeout:                        "5s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "1m", // 1分钟后(最小时间) ,取消注册
	}
	serviceID := fmt.Sprintf("%s-%s-%s", serverName, ipinfo.String(), *port)
	p, _ := strconv.Atoi(*port)
	srv := &api.AgentServiceRegistration{
		ID:      serviceID,
		Name:    serverName,
		Tags:    []string{"aodeibiao"},
		Port:    p,
		Address: "127.0.0.1",
		Check:   check,
	}
	// 2.注册到consul
	cc.Agent().ServiceRegister(srv)
	// 启动服务

	go func() {
		err := grpcServer.Serve(listen)
		if err != nil {
			log.Fatalf("err:%v", err)
			return
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	fmt.Println("wait quit signal....")
	<-quit
	fmt.Println("service out")
	if err := cc.Agent().ServiceDeregister(serviceID); err != nil {
		log.Fatalf("err:%v", err)
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
