package client

import (
	"microservice/timor_grpc"
	"microservice/timor_grpc/pb"

	grpctransport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
)

// Return new timor_grpc service
func New(conn *grpc.ClientConn) timor_grpc.Service {
	var timorEndpoint = grpctransport.NewClient(
		conn, "Timor", "Timor",
		timor_grpc.EncodeGRPCTimorRequest,
		timor_grpc.DecodeGRPCTimorResponse,
		pb.TimorResponse{},
	).Endpoint()
	return timor_grpc.Endpoints{
		TimorEndpoint: timorEndpoint,
	}
}
