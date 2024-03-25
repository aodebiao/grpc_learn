package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"lb_hello_server/pb"
	"log"
	"net"
)

type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, req *pb.HelloRequest) (res *pb.HelloResponse, error error) {
	reply := fmt.Sprintf("hello %s ,你好啊!  [from %s]", req.Name, addr)
	return &pb.HelloResponse{Reply: reply}, nil
}

var port = flag.String("p", "8989", "-p指定端口,默认8989")
var addr string

func main() {
	flag.Parse()
	addr = fmt.Sprintf("127.0.0.1:%s", *port)

	listen, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("listen port %s error:%v\n", *port, err)
		return
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	err = s.Serve(listen)
	if err != nil {
		log.Fatalf("starting server error:%v\n", err)
		return
	}

}
