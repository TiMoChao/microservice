package timor_grpc

import (
	"context"
	"microservice/timor_grpc/pb"
)

//Encode and Decode Timor Request
func EncodeGRPCTimorRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(TimorRequest)
	return &pb.TimorRequest{
		RequestType: req.RequestType,
		Max:         req.Max,
		Min:         req.Min,
	}, nil
}

func DecodeGRPCTimorRequest(ctx context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.TimorRequest)
	return TimorRequest{
		RequestType: req.RequestType,
		Max:         req.Max,
		Min:         req.Min,
	}, nil
}

// Encode and Decode Timor Response
func EncodeGRPCTimorResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(TimorResponse)
	return &pb.TimorResponse{
		Message: resp.Message,
		Err:     resp.Err,
	}, nil
}

func DecodeGRPCTimorResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(*pb.TimorResponse)
	return TimorResponse{
		Message: resp.Message,
		Err:     resp.Err,
	}, nil
}
