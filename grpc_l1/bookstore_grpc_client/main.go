package main

import (
	"bookstore_grpc_client/pb"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func main() {
	// 建立连接
	conn, err := grpc.Dial("127.0.0.1:4567", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("dial err:%v\n", err)
		return
	}
	defer conn.Close()

	c := pb.NewBookStoreClient(conn)

	response, err := c.ListBooks(context.Background(), &pb.ListBookRequest{Shelf: 1})
	if err != nil {
		log.Fatalf("ListBooks err:%v\n", err)
	}
	fmt.Printf("nextPageToken:%v\n", response.NextPageToken)
	books := response.Books
	for i := 0; i < len(books); i++ {
		fmt.Printf("第%d本书,title:%s,author:%s\n", i+1, books[i].Title, books[i].Author)
	}

}
