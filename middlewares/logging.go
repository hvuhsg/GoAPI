package middlewares

import (
	"fmt"
	"log"
	"time"

	"github.com/hvuhsg/goapi/request"
	"github.com/hvuhsg/goapi/responses"
)

type LoggingMiddleware struct{}

func (LoggingMiddleware) Apply(next AppHandler) AppHandler {
	return func(request *request.Request) responses.Response {
		response := next(request)

		scheme := "http"
		if request.HTTPRequest.TLS != nil {
			scheme = "https"
		}
		fullURL := fmt.Sprintf("%s://%s%s", scheme, request.HTTPRequest.Host, request.HTTPRequest.URL.String())

		method := request.HTTPRequest.Method
		path := request.HTTPRequest.URL.Path
		responseSize := len(response.ToBytes())
		remoteAddr := request.HTTPRequest.RemoteAddr
		date := time.Now().Format("2006-01-02 15:04:05")
		userAgent := request.HTTPRequest.UserAgent()
		statusCode := response.StatusCode()

		log.Printf("%s - - [%s] \"%s %s %s\" %d %d \"%s\" \"%s\"\n", remoteAddr, date, method, path, request.HTTPRequest.Proto, statusCode, responseSize, fullURL, userAgent)
		return response
	}
}
