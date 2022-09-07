package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"microservice/timor_grpc"

	grpcClient "microservice/timor_grpc/client"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {
	var (
		grpcAddr = flag.String("addr", ":8081",
			"gRPC address")
	)
	flag.Parse()
	ctx := context.Background()
	conn, err := grpc.Dial(*grpcAddr, grpc.WithInsecure(),
		grpc.WithTimeout(1*time.Second))

	if err != nil {
		log.Fatalln("gRPC dial:", err)
	}
	defer conn.Close()

	timorService := grpcClient.New(conn)
	args := flag.Args()

	var cmd string
	cmd, args = pop(args)

	switch cmd {
	case "timor":
		var requestType, minStr, maxStr string

		requestType, args = pop(args)
		minStr, args = pop(args)
		maxStr, args = pop(args)
		min, _ := strconv.Atoi(minStr)
		max, _ := strconv.Atoi(maxStr)
		timor(ctx, timorService, requestType, min, max)
	default:
		log.Fatalln("unknown command", cmd)
	}
}

// parse command line argument one by one
func pop(s []string) (string, []string) {
	if len(s) == 0 {
		return "", s
	}
	return s[0], s[1:]
}

// call timor service
func timor(ctx context.Context, service timor_grpc.Service, requestType string, min int, max int) {
	mesg, err := service.Timor(ctx, requestType, min, max)
	fmt.Println(mesg)
	os.Exit(1)
	if err != nil {
		log.Fatalln(err.Error())
	}
	fmt.Println(mesg)
}
