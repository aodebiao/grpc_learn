package main

import (
	"addsrv2/pb"
	"flag"
	"fmt"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"net"
	"net/http"
)

var (
	httpAddr = flag.Int("http-addr", 8080, "http端口")
	gRPCAddr = flag.Int("grpc-addr", 8972, "grpc端口")
)

func main() {
	srv := NewService()
	var g errgroup.Group

	g.Go(func() error {
		httpListener, err := net.Listen("tcp", fmt.Sprintf(":%d", *httpAddr))
		if err != nil {
			fmt.Printf("net.Listen %d faield,err:%v\n", *httpAddr, err)
			return err
		}
		defer httpListener.Close()
		httpHandler := NewHttpServer(srv)
		return http.Serve(httpListener, httpHandler)
	})
	g.Go(func() error {

		grpcListener, err := net.Listen("tcp", fmt.Sprintf(":%d", *gRPCAddr))
		if err != nil {
			fmt.Printf("net.listen failed,err:%v\n", err)
			return err
		}
		defer grpcListener.Close()
		s := grpc.NewServer()
		pb.RegisterAddServer(s, NewGRPCServer(srv))
		return s.Serve(grpcListener)

	})
	g.Wait()

}
