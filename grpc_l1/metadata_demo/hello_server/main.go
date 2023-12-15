package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"hello_server_md/pb"
	"log"
	"net"
	"strconv"
	"time"
)

type GreeterServer struct {
	pb.UnimplementedGreeterServer
}

func (s *GreeterServer) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	// trailer是在请求响应后,发送,所以利用defer
	defer func() {
		trailer := metadata.Pairs("timestamp", strconv.Itoa(int(time.Now().Unix())))
		grpc.SetTrailer(ctx, trailer)
	}()
	// 执行业务逻辑之前要check metadata中的是否包含token
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "无效请求")
	}
	// 或者md.Get("token")
	vl := md.Get("token")
	if len(vl) < 1 || vl[0] != "app-test-aodeibiao" {
		return nil, status.Error(codes.Unauthenticated, "无效token")
	}
	//if vl, ok := md["token"]; ok {
	//	if len(vl) > 0 && vl[0] == "app-test-aodeibiao" {
	//		有效请求
	//}
	//}
	reply := "hello " + req.GetName()
	// 发送数据前发送header
	header := metadata.New(map[string]string{
		"location": "chengdu",
	})
	grpc.SendHeader(ctx, header)
	return &pb.HelloResponse{Reply: reply}, nil
}

func main() {
	listen, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Fatalf("listen error:%v", err)
	}
	server := grpc.NewServer()
	pb.RegisterGreeterServer(server, &GreeterServer{})
	err = server.Serve(listen)
	if err != nil {
		log.Fatalf("server start error:%v", err)
	}
}
