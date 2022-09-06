package main

import (
	"flag"
	"fmt"
	"microservice/go_kit_grpc_demo/server/Endpoint"
	"microservice/go_kit_grpc_demo/server/Server"
	"microservice/go_kit_grpc_demo/server/Transport"
	"microservice/go_kit_grpc_demo/server/pb"
	"net"
	"os"
	"os/signal"
	"syscall"

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
	var svc Server.Service
	svc = Server.TimorService{}
	errChan := make(chan error)

	// creating Endpoints struct
	// 译者注: 其实就是另一种形式的路由-控制器映射
	endpoints := Endpoint.Endpoints{
		TimorEndpoint: Endpoint.MakeTimorEndpoint(svc),
	}

	go func() {
		listener, err := net.Listen("tcp", *gRPCAddr)
		if err != nil {
			errChan <- err
			return
		}
		handler := Transport.NewGRPCServer(ctx, endpoints)
		gRPCServer := grpc.NewServer()
		pb.RegisterTimorServer(gRPCServer, handler)
		errChan <- gRPCServer.Serve(listener)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	fmt.Println(<-errChan)
}
