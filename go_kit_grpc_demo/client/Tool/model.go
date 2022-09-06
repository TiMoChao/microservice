package Tool

import (
	"context"
	"microservice/go_kit_grpc_demo/client/pb"
)

//Encode and Decode Timor Request
func EncodeGRPCTimorRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(pb.TimorRequest)
	return &pb.TimorRequest{
		RequestType: req.RequestType,
		Max:         req.Max,
		Min:         req.Min,
	}, nil
}

func DecodeGRPCTimorResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(*pb.TimorResponse)
	return pb.TimorResponse{
		Message: resp.Message,
		Err:     resp.Err,
	}, nil
}
