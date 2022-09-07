package Endpoint

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"microservice/go_kit_grpc_demo/server/Server"
	"microservice/go_kit_grpc_demo/server/pb"
)

type TimorRequest struct {
	RequestType string
	Min         int32
	Max         int32
}

type TimorResponse struct {
	Message string `json:"message"`
	Err     string `json:"err,omitempty"`
}

// 这里仍是传统的MakeXXXEndpoint函数
func MakeTimorEndpoint(svc Server.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(pb.TimorRequest)

		var (
			min, max int
		)

		min = int(req.Min)
		max = int(req.Max)
		txt, err := svc.Timor(ctx, req.RequestType, min, max)

		if err != nil {
			return nil, err
		}

		return TimorResponse{Message: txt}, nil
	}

}

type Endpoints struct {
	TimorEndpoint endpoint.Endpoint
}
