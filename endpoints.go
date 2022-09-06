package microservice

import (
	"context"
	"errors"
	"strings"

	"github.com/go-kit/kit/endpoint"
)

//request
type TimorchaoRequest struct {
	RequestType string
	Min         int
	Max         int
}

//response
type TimorchaoResponse struct {
	Message string `json:"message"`
	Err     error  `json:"err,omitempty"` //omitempty means, if the value is nil then this field won't be displayed
}

//Health Request
type HealthRequest struct {
}

//Health Response
type HealthResponse struct {
	Status bool `json:"status"`
}

var (
	ErrRequestTypeNotFound = errors.New("Request type only valid for word, sentence and paragraph")
)

// endpoints wrapper
type Endpoints struct {
	TimorchaoEndpoint endpoint.Endpoint
	HealthEndpoint    endpoint.Endpoint
}

// creating health endpoint
func MakeHealthEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		status := svc.HealthCheck()
		return HealthResponse{Status: status}, nil
	}
}

// creating timorchao Ipsum Endpoint
func MaketimorchaoEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(TimorchaoRequest)

		var (
			txt      string
			min, max int
		)

		min = req.Min
		max = req.Max

		if strings.EqualFold(req.RequestType, "Word") {
			txt = svc.Word(min, max)
		} else if strings.EqualFold(req.RequestType, "Sentence") {
			txt = svc.Sentence(min, max)
		} else if strings.EqualFold(req.RequestType, "Paragraph") {
			txt = svc.Paragraph(min, max)
		} else {
			return nil, ErrRequestTypeNotFound
		}

		return TimorchaoResponse{Message: txt}, nil
	}

}
