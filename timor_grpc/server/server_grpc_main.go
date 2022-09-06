package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"microservice/timor_grpc"
	"microservice/timor_grpc/pb"

	context "golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {

	var (
		gRPCAddr = flag.String("grpc", ":8081",
			"gRPC listen address")
	)
	flag.Parse()
	ctx := context.Background()

	// init timor service
	var svc timor_grpc.Service
	svc = timor_grpc.TimorService{}
	errChan := make(chan error)

	// creating Endpoints struct
	// 译者注: 其实就是另一种形式的路由-控制器映射
	endpoints := timor_grpc.Endpoints{
		TimorEndpoint: timor_grpc.MakeTimorEndpoint(svc),
	}

	go func() {
		listener, err := net.Listen("tcp", *gRPCAddr)
		if err != nil {
			errChan <- err
			return
		}
		handler := timor_grpc.NewGRPCServer(ctx, endpoints)
		gRPCServer := grpc.NewServer()
		pb.RegisterTimorServer(gRPCServer, handler)
		errChan <- gRPCServer.Serve(listener)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()
	fmt.Println(<-errChan)
}
