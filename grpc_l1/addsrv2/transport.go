package main

import (
	"addsrv2/pb"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	httptransport "github.com/go-kit/kit/transport/http"
	"net/http"
)

// transport

// http json transport
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

func NewHttpServer(svc AddService) http.Handler {
	sumHandler := httptransport.NewServer(makeSumEndpoint(svc), decodeSumRequest, encodeResponse)
	concatHandler := httptransport.NewServer(makeConcatEndpoint(svc), decodeConcatRequest, encodeResponse)
	//http.Handle("/sum", sumHandler)
	//http.Handle("/concat", concatHandler)
	//http.ListenAndServe(":8080", nil)

	r := gin.Default()
	r.POST("/sum", gin.WrapH(sumHandler))
	r.POST("/concat", gin.WrapH(concatHandler))
	return r
}

// grpc  transport

// NewGRPCServer 构造函数
func NewGRPCServer(svc AddService) pb.AddServer {
	return &grpcServer{sum: grpctransport.NewServer(makeSumEndpoint(svc), decodeGRPCSumRequest, encodeGRCPSumResponse),
		concat: grpctransport.NewServer(makeConcatEndpoint(svc), decodeGRPCConcatRequest, encodeGRCPCConcatResponse)}
}

// 网络传输相关的，包括协议（http,grpc,thrift），等
// grpc的请求与响应
// 以下是对请求和响应的编解码
// 将Sum方法的Grpc参数，转为内部的SumRequest
func decodeGRPCSumRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.SumRequest)
	return SumRequest{A: int(req.A), B: int(req.B)}, nil
}

// 将Sum方法的Grpc参数，转为内部的ConcatRequest
func decodeGRPCConcatRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.ConcatRequest)
	return ConcatRequest{A: req.A, B: req.B}, nil
}

func encodeGRCPSumResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(SumResponse)
	return &pb.SumResponse{V: int64(resp.V), Err: resp.Err}, nil
}

func encodeGRCPCConcatResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(ConcatResponse)
	return &pb.ConcatResponse{V: resp.V, Err: resp.Err}, nil

}
