package main

import (
	"bookstore/pb"
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"net/http"
)

func main() {
	// 连接数据库
	db, err := NewDB("bookstore.db")
	if err != nil {
		fmt.Printf("connect to db failed,err:%v\n", err)
		return
	}
	// 创建 server
	srv := server{bs: &bookstore{db: db}}
	listen, err := net.Listen("tcp", ":4567")
	if err != nil {
		log.Fatalf("failed to listen err:%v\n", err)
	}
	s := grpc.NewServer()
	// 注册服务
	pb.RegisterBookStoreServer(s, &srv)
	go func() {
		fmt.Println(s.Serve(listen))
	}()
	// grpc-gateway
	conn, err := grpc.DialContext(context.Background(), "127.0.0.1:4567",
		grpc.WithBlock(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("grpc conn failed,err:%v\n", err)
	}
	gwmux := runtime.NewServeMux()
	pb.RegisterBookStoreHandler(context.Background(), gwmux, conn)
	gwServer := &http.Server{Addr: ":8090", Handler: gwmux}
	fmt.Println("grpc gateway 8090")
	gwServer.ListenAndServe()
}
