package main

import (
	"addsrv2/pb"
	"context"
	"errors"
	"github.com/go-kit/kit/log"
	grpctransport "github.com/go-kit/kit/transport/grpc"
)

// service层

// 所有和业务逻辑相关的逻辑，应该放在这
type AddService interface {
	Sum(ctx context.Context, a, b int) (int, error)
	Concat(ctx context.Context, a, b string) (string, error)
}

var (
	// ErrEmptyString 两个参数都是空字符串
	ErrEmptyString = errors.New("两个参数都是空字符串")
)

// addService 一个AddService接口的具体实现
// 它的内部可以添加各种扩展字段
type addService struct {
}

func NewService() AddService {
	return &addService{}
}

func (as addService) Sum(ctx context.Context, a, b int) (int, error) {
	// 业务逻辑，查库
	return a + b, nil
}

// Concat 拼接两个字符串
func (as addService) Concat(ctx context.Context, a, b string) (string, error) {
	if a == "" && b == "" {
		return "", ErrEmptyString
	}
	return a + b, nil
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

type LogMiddleware struct {
	log  log.Logger
	next AddService
}
