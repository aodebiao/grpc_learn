# 第一个gRPC示例


hello world

## 三个步骤
1.编写protobuf文件 
2.生成代码(服务端和客户端),在项目根目录下执行
```
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/hello.proto

protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative pb/hello.proto
```
3.编写业务逻辑代码

- 客户端引入自动自成文件的包名需要和服务端保持一致,同时目录也需要和option_go_package中的后缀保持一致,不然会报错(粗略描述,请参考metadata_demo中的server和client)
```
 go run main.go 
2023/12/15 19:49:46 client.SayHello error:rpc error: code = Unimplemented desc = unknown service pb.Greeter
exit status 1

```