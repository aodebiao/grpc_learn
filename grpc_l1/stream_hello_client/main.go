package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	"os"
	"stream_hello_client/pb"
	"strings"
	"time"
)

var name = flag.String("name", "熊二", "-name参数指定名称")

func main() {
	flag.Parse()
	dial, err := grpc.Dial("127.0.0.1:8888", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("dial error:%v", err)
	}
	client := pb.NewGreeterClient(dial)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*3)
	defer cancel()
	//serverRespStream(ctx, client)
	//fmt.Println("分割线===============")
	//clientSendStream(ctx, client)
	//fmt.Println("分割线===============")
	runBidiHello(ctx, client)
}

func clientSendStream(ctx context.Context, client pb.GreeterClient) {
	stream, err := client.LotsOfGreetings(ctx)
	if err != nil {
		log.Fatalf("error:%v", err)
	}
	names := []string{"张三", "李四", "王二麻子"}
	for _, name := range names {
		stream.Send(&pb.HelloRequest{Name: name})
	}
	// 流式发送结束,关闭流
	recv, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error:%v", err)
	}
	log.Printf("res:%v\n", recv)
}

func serverRespStream(ctx context.Context, client pb.GreeterClient) {
	replies, err := client.LotsOfReplies(ctx, &pb.HelloRequest{Name: *name})
	if err != nil {
		log.Fatalf("grpc response error:%v", err)
	}
	for {
		recv, err := replies.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("recv error:%v", err)
		}
		fmt.Printf("got reply:%q\n", recv.GetReply())
	}
}

func runBidiHello(ctx context.Context, client pb.GreeterClient) {
	stream, err := client.BidiHello(ctx)
	if err != nil {
		log.Fatalf("client.BidiHello failed,error:%v", err)
	}
	waitc := make(chan struct{})
	go func() {
		for {
			recv, err := stream.Recv()
			if err == io.EOF {
				close(waitc)
				return
			}

			if err != nil {
				log.Fatalf("client.BidiHello stream.recv,error:%v", err)
			}
			fmt.Printf("AI:%s\n", recv.GetReply())
		}
	}()

	reader := bufio.NewReader(os.Stdin)
	for {
		cmd, _ := reader.ReadString('\n')
		cmd = strings.TrimSpace(cmd)
		if len(cmd) == 0 {
			continue
		}
		if strings.ToUpper(cmd) == "QUIT" {
			break
		}
		if err := stream.Send(&pb.HelloRequest{Name: cmd}); err != nil {
			log.Fatalf("client stream.send failed,error:%v", err)
		}
	}
	stream.CloseSend()
	<-waitc
}
