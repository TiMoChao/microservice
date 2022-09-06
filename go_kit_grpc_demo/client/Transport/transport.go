package Transport

import (
	"context"
	"microservice/go_kit_grpc_demo/client/pb"

	grpctransport "github.com/go-kit/kit/transport/grpc"
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
