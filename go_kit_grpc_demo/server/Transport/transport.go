package Transport

import (
	"context"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"microservice/go_kit_grpc_demo/server/Endpoint"
	"microservice/go_kit_grpc_demo/server/Tool"
	"microservice/go_kit_grpc_demo/server/pb"
)

type grpcServer struct {
	timor grpctransport.Handler
}

// implement TimorServer Interface in Timor.pb.go
func (s *grpcServer) Timor(ctx context.Context, r *pb.TimorRequest) (*pb.TimorResponse, error) {
	_, resp, err := s.timor.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.TimorResponse), nil
}

// create new grpc server
func NewGRPCServer(_ context.Context, endpoint Endpoint.Endpoints) pb.TimorServer {
	return &grpcServer{
		timor: grpctransport.NewServer(
			endpoint.TimorEndpoint,
			Tool.DecodeGRPCTimorRequest,
			Tool.EncodeGRPCTimorResponse,
		),
	}
}
