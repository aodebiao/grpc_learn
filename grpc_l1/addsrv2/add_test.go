package main

import (
	"addsrv2/pb"
	"context"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"log"
	"net"
	"testing"
)

// gRpc test

// 编写一个gRPC客户端

// 使用bufconn构建单元测试，避免使用实际商品号启动服务
const bufSize = 1 << 20

var bufListener *bufconn.Listener

func init() {
	bufListener = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	gs := NewGRPCServer(addService{})
	pb.RegisterAddServer(s, gs)
	go func() {
		if err := s.Serve(bufListener); err != nil {
			log.Fatalf("Server exited with error:%v", err)
		}
	}()
}

func bufDialer(ctx context.Context, str string) (net.Conn, error) {
	return bufListener.Dial()

}

// 原始的单元测试
func TestSum(t *testing.T) {
	// 建立连接
	conn, err := grpc.DialContext(context.Background(),
		"bufnet",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(bufDialer))
	if err != nil {
		t.Fail()
	}
	defer conn.Close()
	// 创建客户端
	c := pb.NewAddClient(conn)
	// 调用 grpc
	resp, err := c.Sum(context.Background(), &pb.SumRequest{A: 10, B: 20})
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, resp.V, int64(30))
}

// 原始的单元测试
func TestSum1(t *testing.T) {
	// 建立连接
	conn, err := grpc.DialContext(context.Background(), "127.0.0.1:8090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fail()
	}
	defer conn.Close()
	// 创建客户端
	c := pb.NewAddClient(conn)
	// 调用 grpc
	resp, err := c.Sum(context.Background(), &pb.SumRequest{A: 10, B: 20})
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, resp.V, int64(30))
}
