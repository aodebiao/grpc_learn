package main

import (
	"bookstore/pb"
	"context"
	"testing"
)

// TestServer_ListBooks 测试 server结构体的ListBooks方法
// go test -v
func TestServer_ListBooks(t *testing.T) {
	db, _ := NewDB("test.db")
	s := server{bs: &bookstore{
		db: db,
	}}
	books, err := s.ListBooks(context.Background(), &pb.ListBookRequest{Shelf: 1})
	if err != nil {
		t.Fatalf("ListBooks err:%v\n", err)
	}
	t.Logf("next_page_token:%v\n", books.NextPageToken)
	for i, book := range books.Books {
		t.Logf("第%d本书,title:%s,author:%s,", i, book.Title, book.Author)
	}

}
