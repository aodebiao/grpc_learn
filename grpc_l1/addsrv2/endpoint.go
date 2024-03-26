package main

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

type SumRequest struct {
	A int `json:"a"`
	B int `json:"b"`
}
type SumResponse struct {
	V   int    `json:"v"`
	Err string `json:"err,omitempty"`
}

type ConcatRequest struct {
	A string `json:"a"`
	B string `json:"b"`
}
type ConcatResponse struct {
	V   string `json:"v"`
	Err string `json:"err,omitempty"`
}

// endpoint
// 一个endpoint表示对外提供的方法

// 2将方法转为endpoint
func makeSumEndpoint(s AddService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(SumRequest)
		v, err := s.Sum(ctx, int(req.A), int(req.B))
		if err != nil {
			return SumResponse{V: int(v), Err: err.Error()}, nil
		}
		return SumResponse{V: int(v)}, nil
	}
}
func makeConcatEndpoint(srv AddService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(ConcatRequest)
		v, err := srv.Concat(ctx, req.A, req.B)
		if err != nil {
			return ConcatResponse{V: v, Err: err.Error()}, nil
		}
		return ConcatResponse{V: v}, nil
	}

}
