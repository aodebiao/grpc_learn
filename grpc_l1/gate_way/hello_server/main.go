package main

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"net/http"
	"way_server1/pb"
)

type Server struct {
	pb.UnimplementedGreeterServer
}

func (s Server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	reply := "您好: " + in.Name
	return &pb.HelloResponse{Reply: reply}, nil
}

func main() {
	listen, err := net.Listen("tcp", ":4567")
	if err != nil {
		log.Fatalf("listen failed,error:%v\n", err)
	}
	server := grpc.NewServer()
	pb.RegisterGreeterServer(server, &Server{})
	//err = server.Serve(listen) // 会阻塞
	//if err != nil {
	//	log.Fatalf("serve start failed,error:%v\n", err)
	//}

	go func() {
		log.Fatalln(server.Serve(listen))
	}()

	// 下面是grpc gateway 新加的内容
	// 创建一个连接到我们刚刚启动的grpc 服务器的客户端连接
	// gRPC-GateWay 就是通过它来代理请求,(将http请求转为rpc)请求
	conn, err := grpc.DialContext(context.Background(),
		"127.0.0.1:4567",
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("dial context error :%v\n", err)
	}
	gwmux := runtime.NewServeMux()
	err = pb.RegisterGreeterHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalf("failed to register gateway:%v\n", err)
	}
	gwServer := &http.Server{
		Addr:    ":8090",
		Handler: gwmux,
	}
	log.Println("sering grpc-Gateway on http://127.0.0.1:8090")
	log.Fatalln(gwServer.ListenAndServe())
}
