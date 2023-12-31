package main

import (
	"database/sql"
	"fmt"
	"github.com/iancoleman/strcase"
	fieldmask_utils "github.com/mennanov/fieldmask-utils"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"protobuf_demo/api"
)

// oneofDemo 示例
func oneofDemo() {

	// client
	req := &api.NoticeReaderRequest{
		Msg:       "雾山五行更新了",
		NoticeWay: &api.NoticeReaderRequest_Email{Email: "123@qq.com"},
	}

	req1 := &api.NoticeReaderRequest{
		Msg:       "雾山五行更新了",
		NoticeWay: &api.NoticeReaderRequest_Phone{Phone: "1888888888"},
	}
	// server
	switch v := req.NoticeWay.(type) {
	case *api.NoticeReaderRequest_Email:
		noticeWithEmail(v)
	case *api.NoticeReaderRequest_Phone:
		noticeWithPhone(v)

	}
	fmt.Println("我是分割线...........")
	switch v := req1.NoticeWay.(type) {
	case *api.NoticeReaderRequest_Email:
		noticeWithEmail(v)
	case *api.NoticeReaderRequest_Phone:
		noticeWithPhone(v)

	}

}

type Book struct {
	Price  int64         // 区分默认值和零值
	Price1 sql.NullInt64 // 1:包装(自定义结构体)
	Price2 *int64        // 2:使用指针
}

// 区分默认值和零值
// 1:使用指针
// 2:包装
func foo() {
	var book Book
	if book.Price2 != nil {
		// 默认值
	} else {
		// 零值
	}
}

// wrapValueDemo
// 在grpc中区分零值和默认值,使用google/
func wrapValueDemo() {
	book := api.Book{
		Title: "熊出没",
		Price: &wrapperspb.Int64Value{
			Value: 9999,
		},
		Memo: &wrapperspb.StringValue{Value: "学就完事了"},
	}
	if book.GetPrice() == nil {
		// 零值
	} else {
		// 赋值了
		fmt.Printf("book is sale :%v\n", book.GetPrice().GetValue())
	}
}

func optionalDemo() {
	book := api.Book{
		Title:         "熊出没",
		PriceOptional: proto.Int64(1111),
	}
	if book.PriceOptional == nil {
		fmt.Println("未设置价格")
		return
	}
	fmt.Println("optional value is ", *book.PriceOptional)
}

// fieldMaskDemo 部分更新实例
// 使用fieldMask更新实例
func fieldMaskDemo() {
	// client
	paths := []string{"price", "info.b", "author"} // 更新字段的路径信息
	req := api.UpdateBookRequest{
		Op: "熊二",
		Book: &api.Book{
			Author: "熊大",
			Price:  &wrapperspb.Int64Value{Value: 8800},
			Info:   &api.Book_Info{B: "bbbb"},
		},
		UpdateMask: &fieldmaskpb.FieldMask{Paths: paths},
	}

	// server
	mask, _ := fieldmask_utils.MaskFromProtoFieldMask(req.UpdateMask, strcase.ToCamel)
	fmt.Printf("mask value : %v\n", mask.String())
	var bookDst = make(map[string]any)
	_ = fieldmask_utils.StructToMap(mask, req.Book, bookDst)
	fmt.Printf("bookDst: %#v", bookDst)

}

func main() {
	oneofDemo()
	wrapValueDemo()
	optionalDemo()
	fieldMaskDemo()
}

// 发送通知相关的功能函数
func noticeWithEmail(in *api.NoticeReaderRequest_Email) {
	fmt.Printf("notice reader by email:%v\n", in.Email)
}

func noticeWithPhone(in *api.NoticeReaderRequest_Phone) {
	fmt.Printf("notice reader by phone:%v\n", in.Phone)
}
