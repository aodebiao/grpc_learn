package main

import (
	"addsrv2/pb"
	"context"
	"errors"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
	"net"
	"net/http"
)

type AddService interface {
	Sum(ctx context.Context, a, b int) (int, error)
	Concat(ctx context.Context, a, b string) (string, error)
}

var (
	// ErrEmptyString 两个参数都是空字符串
	ErrEmptyString = errors.New("两个参数都是空字符串")
)

type addService struct{}

func (as addService) Sum(ctx context.Context, a, b int) (int, error) {

	return a + b, nil
}

// Concat 拼接两个字符串
func (as addService) Concat(ctx context.Context, a, b string) (string, error) {
	if a == "" && b == "" {
		return "", ErrEmptyString
	}
	return a + b, nil
}

func ping(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("pong"))
	if err != nil {
		fmt.Println("http write err,", err)
	}
}

func main() {
	srv := addService{}
	gs := NewGRPCServer(srv)

	l, err := net.Listen("tcp", ":8090")
	if err != nil {
		fmt.Printf("net.listen failed,err:%v\n", err)
		return
	}
	s := grpc.NewServer() // grpc server
	pb.RegisterAddServer(s, gs)
	fmt.Println(s.Serve(l))
}

// grpc

type grpcServer struct {
	pb.UnimplementedAddServer
	sum    grpctransport.Handler
	concat grpctransport.Handler
}

func (receiver grpcServer) Sum(ctx context.Context, request *pb.SumRequest) (*pb.SumResponse, error) {
	_, resp, err := receiver.sum.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.SumResponse), nil
}
func (receiver grpcServer) Concat(ctx context.Context, request *pb.ConcatRequest) (*pb.ConcatResponse, error) {
	_, resp, err := receiver.concat.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.ConcatResponse), nil
}

// NewGRPCServer 构造函数
func NewGRPCServer(svc AddService) pb.AddServer {
	return &grpcServer{sum: grpctransport.NewServer(makeSumEndpoint(svc), decodeGRPCSumRequest, encodeGRCPSumResponse),
		concat: grpctransport.NewServer(makeConcatEndpoint(svc), decodeGRPCConcatRequest, encodeGRCPCConcatResponse)}
}

// 2将方法转为endpoint
func makeSumEndpoint(s AddService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(pb.SumRequest)
		v, err := s.Sum(ctx, int(req.A), int(req.B))
		if err != nil {
			return pb.SumResponse{V: int64(v), Err: err.Error()}, nil
		}
		return pb.SumResponse{V: int64(v)}, nil
	}
}
func makeConcatEndpoint(srv AddService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(pb.ConcatRequest)
		v, err := srv.Concat(ctx, req.A, req.B)
		if err != nil {
			return pb.ConcatResponse{V: v, Err: err.Error()}, nil
		}
		return pb.ConcatResponse{V: v}, nil
	}

}

// grpc的请求与响应
// 以下是对请求和响应的编解码
// 将Sum方法的Grpc参数，转为内部的SumRequest
func decodeGRPCSumRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.SumRequest)
	return pb.SumRequest{A: req.A, B: req.B}, nil
}

// 将Sum方法的Grpc参数，转为内部的ConcatRequest
func decodeGRPCConcatRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.ConcatRequest)
	return pb.ConcatRequest{A: req.A, B: req.B}, nil
}

func encodeGRCPSumResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(pb.SumResponse)
	return &pb.SumResponse{V: int64(resp.V), Err: resp.Err}, nil
}

func encodeGRCPCConcatResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(pb.ConcatResponse)
	return &pb.ConcatResponse{V: resp.V, Err: resp.Err}, nil

}
