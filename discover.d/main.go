package main

import (
	"context"
	"io"
	"microservice"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	consulsd "github.com/go-kit/kit/sd/consul"
	"github.com/go-kit/kit/sd/lb"
	ht "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/hashicorp/consul/api"
)

//to execute: go run src/microservice/discover.d/main.go -consul.addr localhost -consul.port 8500
// curl -XPOST -d'{"requestType":"word", "min":10, "max":10}' http://localhost:8080/sd-timor
func main() {

	var (
		consulAddr = flag.String("consul.addr", "", "consul address")
		consulPort = flag.String("consul.port", "", "consul port")
	)
	flag.Parse()

	// Logging domain.
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	// Service discovery domain. In this example we use Consul.
	var client consulsd.Client
	{
		consulConfig := api.DefaultConfig()

		consulConfig.Address = "http://" + *consulAddr + ":" + *consulPort
		consulClient, err := api.NewClient(consulConfig)
		if err != nil {
			logger.Log("err", err)
			os.Exit(1)
		}
		client = consulsd.NewClient(consulClient)
	}

	tags := []string{"microservice"}
	passingOnly := true
	duration := 500 * time.Millisecond
	var timorchaoEndpoint endpoint.Endpoint

	ctx := context.Background()
	r := mux.NewRouter()

	factory := timorchaoFactory(ctx, "GET", "/timorchao")
	serviceName := "microservice"
	instancer := consulsd.NewInstancer(client, logger, serviceName, tags, passingOnly)
	endpointer := sd.NewEndpointer(instancer, factory, logger)
	balancer := lb.NewRoundRobin(endpointer)
	retry := lb.Retry(1, duration, balancer)
	timorchaoEndpoint = retry

	// configure hystrix
	hystrix.ConfigureCommand("timorchao Request", hystrix.CommandConfig{Timeout: 1000})
	timorchaoEndpoint = microservice.Hystrix("timorchao Request",
		"Service currently unavailable", logger)(timorchaoEndpoint)

	// POST /sd-timorchao
	// Payload: {"requestType":"word", "min":10, "max":10}
	r.Methods("POST").Path("/sd-timorchao").Handler(ht.NewServer(
		timorchaoEndpoint,
		decodeConsultimorchaoRequest,
		microservice.EncodeResponse, // use existing encode response since I did not change the logic on response
	))

	// Interrupt handler.
	errc := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-c)
	}()

	// configure the hystrix stream handler
	hystrixStreamHandler := hystrix.NewStreamHandler()
	hystrixStreamHandler.Start()
	go func() {
		errc <- http.ListenAndServe(net.JoinHostPort("", "9000"), hystrixStreamHandler)
	}()

	// HTTP transport.
	go func() {
		logger.Log("transport", "HTTP", "addr", "8080")
		errc <- http.ListenAndServe(":8080", r)
	}()

	// Run!
	logger.Log("exit", <-errc)
}

// factory function to parse URL from Consul to Endpoint
func timorchaoFactory(_ context.Context, method, path string) sd.Factory {
	return func(instance string) (endpoint.Endpoint, io.Closer, error) {
		if !strings.HasPrefix(instance, "http") {
			instance = "http://" + instance
		}

		tgt, err := url.Parse(instance)
		if err != nil {
			return nil, nil, err
		}
		tgt.Path = path

		var (
			enc ht.EncodeRequestFunc
			dec ht.DecodeResponseFunc
		)
		enc, dec = encodetimorchaoRequest, decodetimorchaoResponse

		return ht.NewClient(method, tgt, enc, dec).Endpoint(), nil, nil
	}
}

// decode request from discovery service
// parsing JSON into TimorchaoRequest
func decodeConsultimorchaoRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request microservice.TimorchaoRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

// Encode request form timorchaoRequest into existing timorchao Service
// The encode will translate the timorchaoRequest into /timorchao/{requestType}/{min}/{max}
func encodetimorchaoRequest(_ context.Context, req *http.Request, request interface{}) error {
	lr := request.(microservice.TimorchaoRequest)
	p := "/" + lr.RequestType + "/" + strconv.Itoa(lr.Min) + "/" + strconv.Itoa(lr.Max)
	req.URL.Path += p
	return nil
}

// decode response from timorchao Service
func decodetimorchaoResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response microservice.TimorchaoResponse
	var s map[string]interface{}

	if respCode := resp.StatusCode; respCode >= 400 {
		if err := json.NewDecoder(resp.Body).Decode(&s); err != nil {
			return nil, err
		}
		return nil, errors.New(s["error"].(string) + "\n")
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response, nil
}

// encode error
func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

// errorer is implemented by all concrete response types that may contain
// errors. It allows us to change the HTTP response code without needing to
// trigger an endpoint (transport-level) error.
type errorer interface {
	error() error
}

// encodeResponse is the common method to encode all response types to the
// client.
func EncodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		// Not a Go kit transport error, but a business-logic error.
		// Provide those as HTTP errors.
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}
