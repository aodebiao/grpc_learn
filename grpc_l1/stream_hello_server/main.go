package main

import (
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"stream_hello_server/pb"
	"strings"
)

var _ pb.GreeterServer = (*Greeter)(nil)

type Greeter struct {
	pb.UnimplementedGreeterServer
}

func (g Greeter) BidiHello(stream pb.Greeter_BidiHelloServer) error {
	for {
		recv, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		reply := magic(recv.GetName())
		if err := stream.Send(&pb.HelloResponse{Reply: reply}); err != nil {
			return err
		}
	}
}

func (g Greeter) LotsOfGreetings(stream pb.Greeter_LotsOfGreetingsServer) error {
	reply := "你好"
	for {
		recv, err := stream.Recv()
		if err == io.EOF {
			// 最终统一回复
			return stream.SendAndClose(&pb.HelloResponse{Reply: reply})
		}
		if err != nil {
			return err
		}
		reply += recv.Name
	}
}
func (g Greeter) LotsOfReplies(request *pb.HelloRequest, stream pb.Greeter_LotsOfRepliesServer) error {
	words := []string{
		"你好",
		"hello",
		"お元気ですか",
		"안녕하세요",
	}
	for _, word := range words {
		data := &pb.HelloResponse{Reply: word + request.GetName()}
		if err := stream.Send(data); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	listen, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Fatalf("listen error:%v", err)
	}
	server := grpc.NewServer()
	pb.RegisterGreeterServer(server, &Greeter{})
	err = server.Serve(listen)
	if err != nil {
		log.Fatalf("listen error:%v", err)

	}
}

// magic 一段价值连城的“人工智能”代码
func magic(s string) string {
	s = strings.ReplaceAll(s, "吗", "")
	s = strings.ReplaceAll(s, "嘛", "")
	s = strings.ReplaceAll(s, "吧", "")
	s = strings.ReplaceAll(s, "你", "我")
	s = strings.ReplaceAll(s, "？", "!")
	s = strings.ReplaceAll(s, "?", "!")
	return s
}
