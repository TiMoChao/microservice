package client

import (
	"microservice/go_kit_grpc_demo/client/Endpoint"
	"microservice/go_kit_grpc_demo/client/Tool"
	"microservice/go_kit_grpc_demo/client/pb"
	"microservice/timor_grpc"

	grpctransport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
)

// Return new timor_grpc service
func New(conn *grpc.ClientConn) timor_grpc.Service {
	var timorEndpoint = grpctransport.NewClient(
		conn, "Timor", "Timor",
		Tool.EncodeGRPCTimorRequest,
		Tool.DecodeGRPCTimorResponse,
		pb.TimorResponse{},
	).Endpoint()
	return Endpoint.Endpoints{
		TimorEndpoint: timorEndpoint,
	}
}
