package microservice

import (
	"context"
	"errors"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/metrics"
	"github.com/juju/ratelimit"
)

var ErrLimitExceed = errors.New("Rate Limit Exceed")

func NewTokenBucketLimiter(tb *ratelimit.Bucket) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			if tb.TakeAvailable(1) == 0 {
				return nil, ErrLimitExceed
			}
			return next(ctx, request)
		}
	}
}

type metricsMiddleware struct {
	Service
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
}

// metrics function
func Metrics(requestCount metrics.Counter,
	requestLatency metrics.Histogram) ServiceMiddleware {
	return func(next Service) Service {
		return metricsMiddleware{
			next,
			requestCount,
			requestLatency,
		}
	}
}

// Implement service functions and add label method for our metrics
func (mw metricsMiddleware) Word(min, max int) (output string) {
	defer func(begin time.Time) {
		lvs := []string{"method", "Word"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	output = mw.Service.Word(min, max)
	return
}

// and the rest for Sentence and Paragraph

// implement metrics feature in HealthCheck function
func (mw metricsMiddleware) HealthCheck() (output bool) {
	defer func(begin time.Time) {
		lvs := []string{"method", "HealthCheck"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	output = mw.Service.HealthCheck()
	return
}
