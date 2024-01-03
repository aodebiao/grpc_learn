package main

import (
	"bookstore/pb"
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"net/http"
	"strings"
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
	// 同一个端口分别处理grpc和http
	// 1 创建 grpc-gateway mux
	gwmux := runtime.NewServeMux()
	dops := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	if err := pb.RegisterBookStoreHandlerFromEndpoint(context.Background(), gwmux, "127.0.0.1:4567", dops); err != nil {
		log.Fatalf("RegisterBookStoreHandlerFromEndpoint failed,err:%v\n", err)
		return
	}
	// 2 新建 http mux
	mux := http.NewServeMux()
	mux.Handle("/", gwmux)

	// 3 定义http server
	gwServer := &http.Server{Addr: "127.0.0.1:4567", Handler: grpcHandlerFunc(s, mux)}
	// 4 启动
	fmt.Println("staring server on 127.0.0.1:4567")
	gwServer.Serve(listen)
}

// grpcHandlerFunc 将gRPC请求和HTTP请求分别调用不同的handler处理
func grpcHandlerFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	return h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			otherHandler.ServeHTTP(w, r)
		}
	}), &http2.Server{})
}

// httpAndGrpcNotPort 同时提供grpc服务和http服务,但是不在同一个端口
func httpAndGrpcNotPort() {
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
