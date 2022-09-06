package timor_grpc

import (
	"context"
	"microservice/timor_grpc/pb"

	grpctransport "github.com/go-kit/kit/transport/grpc"
)

type grpcServer struct {
	timor grpctransport.Handler
}

// implement Timorerver Interface in timor.pb.go
func (s *grpcServer) Timor(ctx context.Context, r *pb.TimorRequest) (*pb.TimorResponse, error) {
	_, resp, err := s.timor.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.TimorResponse), nil
}

// create new grpc server
func NewGRPCServer(_ context.Context, endpoint Endpoints) pb.TimorServer {
	return &grpcServer{
		timor: grpctransport.NewServer(
			endpoint.TimorEndpoint,
			DecodeGRPCTimorRequest,
			EncodeGRPCTimorResponse,
		),
	}
}
