package timor_grpc

import (
	"context"
	"errors"

	"github.com/go-kit/kit/endpoint"
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
func MakeTimorEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(TimorRequest)

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

// 这里Endpoints也实现了业务层的`Service`接口, 之后将由grpc client直接调用
// 译者注: 其实这里已经不算是endpoint部分了, 而是类似于示例04的transport部分的handler,
// 把不同的grpc请求都交给内部的handler, 调用业务层对象完成逻辑.
func (e Endpoints) Timor(ctx context.Context, requestType string, min, max int) (string, error) {
	req := TimorRequest{
		RequestType: requestType,
		Min:         int32(min),
		Max:         int32(max),
	}
	resp, err := e.TimorEndpoint(ctx, req)
	if err != nil {
		return "", err
	}
	timorResp := resp.(TimorResponse)
	if timorResp.Err != "" {
		return "", errors.New(timorResp.Err)
	}
	return timorResp.Message, nil
}
