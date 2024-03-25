package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"hello_server_interceptor/pb"
	"log"
	"net"
	"strings"
	"time"
)

type server struct {
	pb.UnimplementedGreeterServer
}

type wrapperStream struct {
	grpc.ServerStream
}

func (w *wrapperStream) RecvMsg(m any) error {
	log.Printf("receive a message (Type:%T) at %s\n", m, time.Now().Format(time.RFC3339))
	return w.ServerStream.RecvMsg(m)
}
func (w *wrapperStream) SendMsg(m any) error {
	log.Printf("send a message (Type:%T) at %s\n", m, time.Now().Format(time.RFC3339))
	return w.ServerStream.SendMsg(m)
}

func newWrapperStream(s grpc.ServerStream) grpc.ServerStream {
	return &wrapperStream{s}
}

func valid(authorization []string) bool {
	if len(authorization) < 1 {
		return false
	}
	token := strings.TrimPrefix(authorization[0], "Bearer ")
	return token == "some-secret-token"
}

func unaryInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "missing metadata")
	}
	if !valid(md["authorization"]) {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token")
	}
	m, err := handler(ctx, req)
	if err != nil {
		log.Printf("RPC failed with error %v\n", err)
	}
	return m, err
}

func streamInterceptor(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	md, ok := metadata.FromIncomingContext(ss.Context())
	if !ok {
		return status.Errorf(codes.InvalidArgument, "missing metadata")
	}
	if !valid(md["authorization"]) {
		return status.Errorf(codes.Unauthenticated, "invalid token")
	}
	err := handler(srv, newWrapperStream(ss))
	if err != nil {
		fmt.Printf("RPC failed with error %v\n", err)
	}
	return err
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	reply := "hello " + in.Name
	return &pb.HelloResponse{Reply: reply}, nil
}

func main() {
	listen, err := net.Listen("tcp", ":6789")
	if err != nil {
		log.Fatalf("failed to listen,error:%v\n", err)
	}
	creds, err := credentials.NewServerTLSFromFile("./server.crt", "server.key")
	if err != nil {
		log.Fatalf("tls from file failed,err:%v\n", err)
		return
	}
	s := grpc.NewServer(grpc.Creds(creds),
		grpc.UnaryInterceptor(unaryInterceptor),
		grpc.StreamInterceptor(streamInterceptor))
	pb.RegisterGreeterServer(s, &server{})
	err = s.Serve(listen)
	if err != nil {
		log.Fatalf("failed to server,err:%v\n", err)
		return
	}

}
