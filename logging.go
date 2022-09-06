package microservice

import (
	"time"

	"github.com/go-kit/kit/log"
)

// Make a new type and wrap into Service interface
// Add logger property to this type
type loggingMiddleware struct {
	Service
	logger log.Logger
}

// implement function to return ServiceMiddleware
func LoggingMiddleware(logger log.Logger) ServiceMiddleware {
	return func(next Service) Service {
		return loggingMiddleware{next, logger}
	}
}

// Implement Service Interface for LoggingMiddleware
func (mw loggingMiddleware) Word(min, max int) (output string) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"function", "Word",
			"min", min,
			"max", max,
			"result", output,
			"took", time.Since(begin),
		)
	}(time.Now())
	output = mw.Service.Word(min, max)
	return
}

// Implement Service Interface for LoggingMiddleware
func (mw loggingMiddleware) Sentence(min, max int) (output string) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"function", "Word",
			"min", min,
			"max", max,
			"result", output,
			"took", time.Since(begin),
		)
	}(time.Now())
	output = mw.Service.Sentence(min, max)
	return
}

// implement logging feature in HealthCheck function
func (mw loggingMiddleware) HealthCheck() (output bool) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"function", "HealthCheck",
			"result", output,
			"took", time.Since(begin),
		)
	}(time.Now())
	output = mw.Service.HealthCheck()
	return
}

// and the rest for sentence and paragraph
