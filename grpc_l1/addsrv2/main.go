package main

import (
	"addsrv2/pb"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-kit/kit/log"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"os"
)

var (
	httpAddr = flag.Int("http-addr", 8080, "http端口")
	gRPCAddr = flag.Int("grpc-addr", 8972, "grpc端口")
)

func main() {
	srv := NewService()
	var g errgroup.Group

	// metrics
	// instrumentation
	fieldKeys := []string{"method", "error"}
	requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "add_srv",
		Subsystem: "string_service",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)
	requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "add_srv",
		Subsystem: "string_service",
		Name:      "request_latency_microseconds",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)
	countResult := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "add_srv",
		Subsystem: "string_service",
		Name:      "count_result",
		Help:      "The result of each count method.",
	}, []string{}) // no fields here

	srv = instrumentingMiddleware{
		requestCount:   requestCount,
		requestLatency: requestLatency,
		countResult:    countResult,
		next:           srv,
	}

	g.Go(func() error {
		httpListener, err := net.Listen("tcp", fmt.Sprintf(":%d", *httpAddr))
		if err != nil {
			fmt.Printf("net.Listen %d faield,err:%v\n", *httpAddr, err)
			return err
		}
		defer httpListener.Close()
		logger := log.NewLogfmtLogger(os.Stdout)
		httpHandler := NewHttpServer(srv, logger)
		httpHandler.(*gin.Engine).GET("/metrics", gin.WrapH(promhttp.Handler()))

		// http版本
		//http.Handle("/metrics", promhttp.Handler())
		return http.Serve(httpListener, httpHandler)
	})
	g.Go(func() error {

		grpcListener, err := net.Listen("tcp", fmt.Sprintf(":%d", *gRPCAddr))
		if err != nil {
			fmt.Printf("net.listen failed,err:%v\n", err)
			return err
		}
		defer grpcListener.Close()
		s := grpc.NewServer()
		pb.RegisterAddServer(s, NewGRPCServer(srv))
		return s.Serve(grpcListener)

	})
	g.Wait()

}
