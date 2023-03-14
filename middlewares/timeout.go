package middlewares

import (
	"context"
	"net/http"
	"time"

	"github.com/hvuhsg/goapi/request"
	"github.com/hvuhsg/goapi/responses"
)

type TimeoutMiddleware struct {
	Timeout time.Duration
}

func (tm TimeoutMiddleware) Apply(next AppHandler) AppHandler {
	return func(request *request.Request) responses.Response {
		// Create a context with a timeout
		ctx, cancel := context.WithTimeout(context.Background(), tm.Timeout)
		defer cancel()

		// Add context to the request
		request.HTTPRequest = request.HTTPRequest.WithContext(ctx)

		// Use a goroutine to execute the request and wait for it to complete or timeout
		ch := make(chan responses.Response, 1)
		go func() {
			ch <- next(request)
		}()

		select {
		case response := <-ch:
			return response
		case <-ctx.Done():
			return responses.NewHTMLResponse("Request timed out", http.StatusGatewayTimeout)
		}
	}
}
