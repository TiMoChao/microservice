package Tool

import (
	"context"
	"microservice/go_kit_grpc_demo/server/Endpoint"
	"microservice/go_kit_grpc_demo/server/pb"
)

func DecodeGRPCTimorRequest(ctx context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.TimorRequest)
	return pb.TimorRequest{
		RequestType: req.RequestType,
		Max:         req.Max,
		Min:         req.Min,
	}, nil
}

// Encode and Decode Timor Response
func EncodeGRPCTimorResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(Endpoint.TimorResponse)
	return &pb.TimorResponse{
		Message: resp.Message,
		Err:     resp.Err,
	}, nil
}
