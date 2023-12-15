package main

import (
	"context"
	"fmt"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
	"hello_server/pb"
	"log"
	"net"
	"sync"
)

// grpc server
// 限制一个用户只能调用一次接口
type server struct {
	pb.UnimplementedGreeterServer
	mu    sync.Mutex     // 保证count并发安全
	count map[string]int // 非并发安全
}

// SayHello 需要实现的方法
// 这个方法是我们对外提供的服务
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.count[in.Name]++
	if s.count[in.Name] > 1 {
		//st := status.Error(codes.ResourceExhausted, "request limit")
		//return nil, st

		// 添加错误详情信息
		st := status.New(codes.ResourceExhausted, "request limit")
		// 需要接收新的返回值
		ds, err := st.WithDetails(
			&errdetails.QuotaFailure{
				Violations: []*errdetails.QuotaFailure_Violation{
					{
						Subject:     fmt.Sprintf("name:%s", in.Name),
						Description: "每个name只能调用一次SayHello",
					},
				}},
		)
		// 构造ds失败,返回原error
		if err != nil {
			return nil, st.Err()
		}
		return nil, ds.Err()
	}
	reply := "hello  " + in.GetName()
	return &pb.HelloResponse{Reply: reply}, nil
}

func main() {
	// 启动服务
	l, err := net.Listen("tcp", ":8972")
	if err != nil {
		log.Fatal("failed to listen,err:", err)
	}
	creds, err := credentials.NewServerTLSFromFile("certs/server.crt", "certs/server.key")
	if err != nil {
		log.Fatalf("load certs error:%v", err)
	}

	s := grpc.NewServer(grpc.Creds(creds)) // 创建 grpc服务
	// 注册服务
	pb.RegisterGreeterServer(s, &server{count: make(map[string]int)})
	// 启动服务
	err = s.Serve(l)
	if err != nil {
		log.Fatal("failed to serve,err:", err)
	}
}

// 不安全的连接
func insecureConn() {
	// 启动服务
	l, err := net.Listen("tcp", ":8972")
	if err != nil {
		log.Fatal("failed to listen,err:", err)
	}
	s := grpc.NewServer() // 创建 grpc服务
	// 注册服务
	pb.RegisterGreeterServer(s, &server{count: make(map[string]int)})
	// 启动服务
	err = s.Serve(l)
	if err != nil {
		log.Fatal("failed to serve,err:", err)
	}
}
