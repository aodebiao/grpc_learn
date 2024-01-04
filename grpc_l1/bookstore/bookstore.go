package main

import (
	"bookstore/pb"
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
	"strconv"
	"time"
)

// grpc相关

const defaultPageSize = 2
const defaultCursor = "0" // 默认每页显示数量

type server struct {
	pb.UnimplementedBookStoreServer
	bs *bookstore
}

// ListShelves 列出书架
func (s *server) ListShelves(ctx context.Context, in *emptypb.Empty) (*pb.ListShelvesResponse, error) {
	sl, err := s.bs.ListShelves(ctx)
	if errors.Is(err, gorm.ErrEmptySlice) {
		return nil, nil
	}
	if err != nil {
		return nil, status.Error(codes.Internal, "query failed")
	}
	// 封装返回数据
	nsl := make([]*pb.Shelf, 0, len(*sl))
	for _, s := range *sl {
		nsl = append(nsl, &pb.Shelf{Id: s.ID,
			Theme: s.Theme,
			Size:  s.Size,
		})
	}
	return &pb.ListShelvesResponse{
		Shelves: nsl,
	}, nil
}

// CreateShelf 创建书架
func (s *server) CreateShelf(ctx context.Context, in *pb.CreateShelfRequest) (*pb.Shelf, error) {
	// 参数检查
	if len(in.Shelf.Theme) == 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid theme")
	}
	data := Shelf{
		Theme: in.Shelf.Theme,
		Size:  in.Shelf.Size,
	}
	ns, err := s.bs.CreateShelf(ctx, data)
	if err != nil {
		return nil, status.Error(codes.Internal, "create shelf failed")
	}

	return &pb.Shelf{
		Id:    ns.ID,
		Theme: ns.Theme,
		Size:  ns.Size,
	}, nil

}

// GetShelf 获取书架
func (s *server) GetShelf(ctx context.Context, in *pb.GetShelfRequest) (*pb.Shelf, error) {
	if in.Shelf <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid shelf id")
	}
	shelf, err := s.bs.GetShelf(ctx, in.Shelf)
	if err != nil {
		return nil, status.Error(codes.Internal, "query failed")
	}

	return &pb.Shelf{Id: shelf.ID, Theme: shelf.Theme, Size: shelf.Size}, nil

}

// DeleteShelf 删除书架
func (s *server) DeleteShelf(ctx context.Context, in *pb.DeleteShelfRequest) (*emptypb.Empty, error) {
	if in.Shelf <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid shelf id")
	}
	err := s.bs.DeleteShelf(ctx, in.Shelf)
	if err != nil {
		return nil, status.Error(codes.Internal, "delete failed")
	}
	return &emptypb.Empty{}, nil

}

// ListBooks 列出书架所有图书
func (s *server) ListBooks(ctx context.Context, in *pb.ListBookRequest) (*pb.ListBookResponse, error) {
	// 参数check
	if in.Shelf <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid shelf id")
	}
	// 没有分页token
	var (
		cursor   = defaultCursor
		pageSize = defaultPageSize
	)
	//if pageToken := in.PageToken; pageToken == "" {
	//	// 默认第一页
	//
	//} else {
	if len(in.PageToken) > 0 {
		pageInfo := Token(in.PageToken).Decode()
		// 判断解析结果是否有效
		if pageInfo.InValid() {
			return nil, status.Error(codes.InvalidArgument, "invalid page_token")
		}
		cursor = pageInfo.NextID
		pageSize = int(pageInfo.PageSize)
	}
	//  基于游标实现分页,每次查询的时候多查询一条,判断是否有下一页
	bookList, err := s.bs.GetBookListByShelfID(ctx, in.Shelf, cursor, pageSize+1)
	if err != nil {
		fmt.Printf("GetBookListByShelfID failed,err:%v\n", err)
		return nil, status.Error(codes.Internal, "query failed")
	}
	var (
		hasNextPage   bool
		nextPageToken string
		realSize      int = len(*bookList)
	)
	if len(*bookList) > pageSize {
		hasNextPage = true
		realSize = pageSize
	}

	// 封装返回数据
	res := make([]*pb.Book, 0, len(*bookList))
	//for _, v := range *bookList {
	for i := 0; i < realSize; i++ {
		res = append(res, &pb.Book{Id: (*bookList)[i].ID, Author: (*bookList)[i].Author, Title: (*bookList)[i].Title})
	}
	// 如果有下一页,生成page_token
	if hasNextPage {
		nextPageInfo := Page{PageSize: int64(pageSize), NextID: strconv.FormatInt((*bookList)[realSize-1].ID, 10), NextTimeAtUTC: time.Now().Unix()}
		nextPageToken = string(nextPageInfo.Encode())
	}
	return &pb.ListBookResponse{Books: res, NextPageToken: nextPageToken}, nil
}

// todo实现其它接口
func (s *server) CreateBook(ctx context.Context, in *pb.CreateBookRequest) (*pb.Book, error) {

	return nil, nil

}
func (s *server) GetBook(ctx context.Context, in *pb.GetBookRequest) (*pb.Book, error) {

	return nil, nil

}
func (s *server) DeleteBook(ctx context.Context, in *pb.DeleteBookRequest) (*emptypb.Empty, error) {
	return nil, nil
}
