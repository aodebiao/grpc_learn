package main

import (
	"context"
	"flag"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
	"hello_client_interceptor/pb"
	"log"
	"time"
)

var name = flag.String("name", "熊二", "通过-name告诉server你是谁")

type wrapperStream struct {
	grpc.ClientStream
}

func (w *wrapperStream) RecvMsg(m any) error {
	log.Printf("Receive a message (Type: %T) at %v\n", m, time.Now().Format(time.RFC3339))
	return w.ClientStream.RecvMsg(m)
}

func (w *wrapperStream) SendMsg(m any) error {
	log.Printf("Send a message (Type: %T) at %v\n", m, time.Now().Format(time.RFC3339))
	return w.ClientStream.SendMsg(m)
}
func NewWrapperStream(s grpc.ClientStream) grpc.ClientStream {
	return &wrapperStream{s}
}

func unaryInterceptor(ctx context.Context, method string, req, reply any,
	cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	var credConfigured bool
	for _, o := range opts {
		_, ok := o.(grpc.PerRPCCredsCallOption)
		if ok {
			credConfigured = true
			break
		}
	}
	if !credConfigured {
		opts = append(opts, grpc.PerRPCCredentials(oauth.NewOauthAccess(&oauth2.Token{AccessToken: "some-secret-token"})))
	}
	start := time.Now()
	err := invoker(ctx, method, req, reply, cc, opts...)
	//end := time.Now()
	log.Printf("RPC: %s,duration:%s,err:%v\n", method, time.Since(start), err)
	return err
}

func streamInterceptor(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn,
	method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	var credsConfigured bool
	for _, o := range opts {
		_, ok := o.(*grpc.PerRPCCredsCallOption)
		if ok {
			credsConfigured = true
			break
		}
	}
	if !credsConfigured {
		opts = append(opts, grpc.PerRPCCredentials(oauth.NewOauthAccess(&oauth2.Token{AccessToken: "some-secret-token"})))
	}
	s, err := streamer(ctx, desc, cc, method, opts...)
	if err != nil {
		return nil, err
	}
	return NewWrapperStream(s), nil
}

func main() {
	flag.Parse()
	creds, err := credentials.NewClientTLSFromFile("./server.crt", "")
	if err != nil {
		log.Fatalf("failed to create credentials:%v", err)
	}
	conn, err := grpc.Dial("127.0.0.1:6789",
		grpc.WithTransportCredentials(creds),
		grpc.WithUnaryInterceptor(unaryInterceptor),
		grpc.WithStreamInterceptor(streamInterceptor))
	if err != nil {
		log.Fatalf("grpc.Dial failed,error:%v\n", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	response, err := c.SayHello(ctx, &pb.HelloRequest{Name: *name})
	if err != nil {
		log.Fatalf("c.SayHello failed,err:%v\n", err)
	}
	log.Printf("resp:%v", response.GetReply())
}
