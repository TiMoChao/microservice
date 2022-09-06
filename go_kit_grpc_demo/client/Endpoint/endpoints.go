package Endpoint

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
	TimorResp := resp.(TimorResponse)
	if TimorResp.Err != "" {
		return "", errors.New(TimorResp.Err)
	}
	return TimorResp.Message, nil
}
