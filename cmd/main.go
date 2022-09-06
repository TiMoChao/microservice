package main

import (
	"context"
	"flag"
	"fmt"
	"microservice"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-kit/kit/log"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/juju/ratelimit"
	"github.com/prometheus/client_golang/prometheus"
)

func main() {
	ctx := context.Background()
	errChan := make(chan error)

	// Logging domain.
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "ts", "caller", log.DefaultCaller)
	}

	//declare metrics
	fieldKeys := []string{"method"}
	requestCount := kitprometheus.NewCounterFrom(prometheus.CounterOpts{
		Namespace: "ru_rocker",
		Subsystem: "timorchao_service",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)
	requestLatency := kitprometheus.NewSummaryFrom(prometheus.SummaryOpts{
		Namespace: "ru_rocker",
		Subsystem: "timorchao_service",
		Name:      "request_latency_microseconds",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)

	var svc microservice.Service
	svc = microservice.TimorchaoService{}
	svc = microservice.LoggingMiddleware(logger)(svc)
	svc = microservice.Metrics(requestCount, requestLatency)(svc)

	rlbucket := ratelimit.NewBucket(5*time.Second, 3)
	e1 := microservice.MaketimorchaoEndpoint(svc)
	e1 = microservice.NewTokenBucketLimiter(rlbucket)(e1)

	e2 := microservice.MakeHealthEndpoint(svc)
	e2 = microservice.NewTokenBucketLimiter(rlbucket)(e2)

	endpoint := microservice.Endpoints{
		TimorchaoEndpoint: e1,
		HealthEndpoint:    e2,
	}

	var (
		consulAddr    = flag.String("consul.addr", "", "consul address")
		consulPort    = flag.String("consul.port", "", "consul port")
		advertiseAddr = flag.String("advertise.addr", "", "advertise address")
		advertisePort = flag.String("advertise.port", "", "advertise port")
	)
	flag.Parse()

	// Register Service to Consul
	registrar := microservice.Register(*consulAddr,
		*consulPort,
		*advertiseAddr,
		*advertisePort)

	r := microservice.MakeHttpHandler(ctx, endpoint, logger)

	// HTTP transport
	go func() {
		fmt.Println("Starting server at port 7003")

		// register service
		// 高亮行
		registrar.Register()
		handler := r
		errChan <- http.ListenAndServe(":7010", handler)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()
	fmt.Println(<-errChan)

	// deregister service
	// 高亮行
	registrar.Deregister()
}
