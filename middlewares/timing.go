package middlewares

import (
	"log"
	"time"

	"github.com/hvuhsg/goapi/request"
	"github.com/hvuhsg/goapi/responses"
)

type TimingMiddleware struct{}

func (TimingMiddleware) Apply(next AppHandler) AppHandler {
	return func(request *request.Request) responses.Response {
		startTime := time.Now()
		response := next(request)
		duration := time.Since(startTime)

		log.Printf("goapi.Request took %s\n", duration.String())

		return response
	}
}
