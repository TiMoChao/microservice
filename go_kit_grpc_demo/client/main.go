package main

import (
	"flag"
	"fmt"
	"log"
	"microservice/go_kit_grpc_demo/client/pb"
	"strconv"
	"time"

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

	timorService := pb.NewTimorClient(conn)
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

		stringReq := &pb.TimorRequest{
			RequestType: requestType,
			Min:         int32(min),
			Max:         int32(max),
		}
		reply, err := timorService.Timor(ctx, stringReq)
		if err != nil {
			log.Fatalln(err.Error())
		}
		fmt.Println(reply.Message)
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
