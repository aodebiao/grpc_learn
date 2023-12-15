package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"log"
	"test/proto"
	"time"
)

// grpc客户端
// 调用server端的 SayHello方法
var name = flag.String("name", "aodebiao", "通过-name告诉server你是谁")

func main() {
	flag.Parse()

	// 连接

	// 加载证书,第二个参数如果不为空的话,就不能随便写,需要是生成证书时,指定的SANs中的一个
	creds, err := credentials.NewClientTLSFromFile("certs/server.crt", "127.0.0.1")
	if err != nil {
		log.Fatalf("client load certs failed,error:%v", err)
	}
	conn, err := grpc.Dial("127.0.0.1:8972",
		grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("grpc.Dial failed,err:%v", err)
	}

	defer conn.Close()
	c := proto.NewGreeterClient(conn) // 使用生成的代码
	// 调用RPC方法
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	resp, err := c.SayHello(ctx, &proto.HelloRequest{Name: *name})
	if err != nil {
		// 收到带detail的error
		s := status.Convert(err)
		for _, d := range s.Details() {
			switch info := d.(type) {
			case *errdetails.QuotaFailure:
				fmt.Printf("QuotaFailure:%s\n", info)
			default:
				fmt.Printf("unexpected type:%v\n", err)
			}
		}
		log.Printf("c.SayHello failed,err:%v", err)
		return
	}
	// 拿到结果
	log.Printf("resp value is :%v", resp.Name)
}

func insecureConn() {
	// 连接
	conn, err := grpc.Dial("127.0.0.1:8972", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("grpc.Dial failed,err:%v", err)
	}

	defer conn.Close()
	c := proto.NewGreeterClient(conn) // 使用生成的代码
	// 调用RPC方法
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	resp, err := c.SayHello(ctx, &proto.HelloRequest{Name: *name})
	if err != nil {
		// 收到带detail的error
		s := status.Convert(err)
		for _, d := range s.Details() {
			switch info := d.(type) {
			case *errdetails.QuotaFailure:
				fmt.Printf("QuotaFailure:%s\n", info)
			default:
				fmt.Printf("unexpected type:%v\n", err)
			}
		}
		log.Printf("c.SayHello failed,err:%v", err)
		return
	}
	// 拿到结果
	log.Printf("resp value is :%v", resp.Name)
}
