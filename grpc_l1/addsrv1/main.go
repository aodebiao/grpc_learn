package main

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"net/http"
)

// go-kit addService demo1

// 1.1业务逻辑抽象为接口

type AddService interface {
	Sum(ctx context.Context, a, b int) (int, error)
	Concat(ctx context.Context, a, b string) (string, error)
}

// 1.2实现接口
type addService struct{}

var (
	// ErrEmptyString 两个参数都是空字符串
	ErrEmptyString = errors.New("两个参数都是空字符串")
)

// Sum 返回两个数之和
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

// 1.3请求和响应

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

// 2将方法转为endpoint
func makeSumEndpoint(s AddService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(SumRequest)
		v, err := s.Sum(ctx, req.A, req.B)
		if err != nil {
			return SumResponse{V: v, Err: err.Error()}, nil
		}
		return SumResponse{V: v}, nil
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

// 3.transport

// decode
// 请求来了之后,根据协议(http,http2)和编码(json,pb,thrift,jsonp),去解析数据
func decodeSumRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var request SumRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeConcatRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var request ConcatRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

// encode
// 把响应数据按协议和编码返回
// w:代表响应的网络句柄
// response:代表业务层返回的响应数据
func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func main() {
	srv := addService{}
	sumHandler := httptransport.NewServer(makeSumEndpoint(srv), decodeSumRequest, encodeResponse)
	concatHandler := httptransport.NewServer(makeConcatEndpoint(srv), decodeConcatRequest, encodeResponse)
	http.Handle("/sum", sumHandler)
	http.Handle("/concat", concatHandler)
	http.ListenAndServe(":8080", nil)
}
