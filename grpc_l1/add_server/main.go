package main

import (
	"add_server/proto"
	"context"
	"google.golang.org/grpc"
	"log"
	"net"
)

type AddServer struct {
	proto.UnimplementedAddServerServer
}

func (a AddServer) Add(ctx context.Context, in *proto.AddRequest) (*proto.AddResponse, error) {
	//TODO implement me
	result := in.X + in.Y
	return &proto.AddResponse{Result: result}, nil
}

func main() {
	listen, err := net.Listen("tcp", "127.0.0.1:8888")
	if err != nil {
		log.Fatalf("listen failed,err:%v", err)
	}
	server := grpc.NewServer()
	proto.RegisterAddServerServer(server, &AddServer{})
	err = server.Serve(listen)
	if err != nil {
		log.Fatalf("failed to serve,err:%v", err)
	}
}
