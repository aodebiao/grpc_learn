## proto练习


```
--proto_path 参数指定的protobuf文件的位置,结合最后的book/price.proto
可以找到proto文件并生成代码
--proto_path参数有别名 -I


 protoc --proto_path=proto \
> --go_out=proto \
> --go_opt=paths=source_relative \
> book/price.proto

```