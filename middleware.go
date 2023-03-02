package goapi

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type middleware interface {
	Apply(AppHandler) AppHandler
}

type SimpleLoggingMiddleware struct{}

func (SimpleLoggingMiddleware) Apply(next AppHandler) AppHandler {
	return func(request *Request) Response {
		response := next(request)

		scheme := "http"
		if request.HTTPRequest.TLS != nil {
			scheme = "https"
		}
		fullURL := fmt.Sprintf("%s://%s%s", scheme, request.HTTPRequest.Host, request.HTTPRequest.URL.String())

		method := request.HTTPRequest.Method
		path := request.HTTPRequest.URL.Path
		responseSize := len(response.toBytes())
		remoteAddr := request.HTTPRequest.RemoteAddr
		date := time.Now().Format("2006-01-02 15:04:05")
		userAgent := request.HTTPRequest.UserAgent()
		statusCode := response.statusCode()

		fmt.Printf("%s - - [%s] \"%s %s %s\" %d %d \"%s\" \"%s\"\n", remoteAddr, date, method, path, request.HTTPRequest.Proto, statusCode, responseSize, fullURL, userAgent)
		return response
	}
}

type TimeoutMiddleware struct {
	Timeout time.Duration
}

func (tm TimeoutMiddleware) Apply(next AppHandler) AppHandler {
	return func(request *Request) Response {
		// Create a context with a timeout
		ctx, cancel := context.WithTimeout(context.Background(), tm.Timeout)
		defer cancel()

		// Add context to the request
		request.HTTPRequest = request.HTTPRequest.WithContext(ctx)

		// Use a goroutine to execute the request and wait for it to complete or timeout
		ch := make(chan Response, 1)
		go func() {
			ch <- next(request)
		}()

		select {
		case response := <-ch:
			return response
		case <-ctx.Done():
			return HtmlResponse{Content: "Request timed out", Code: http.StatusGatewayTimeout}
		}
	}
}
