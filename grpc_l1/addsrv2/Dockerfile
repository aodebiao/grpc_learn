FROM golang:1.22-alpine as builder

WORKDIR /app

# 把当前Dockerfile文件所在目录下的文件，拷贝到WORKDIR工作目录中
COPY . .

RUN go mod tidy
RUN go build -o /simple_server

FROM scratch
COPY --from=builder /simple_server /

CMD ["/simple_server"]