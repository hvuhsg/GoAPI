package goapi

import (
	"net/http"

	"github.com/hvuhsg/goapi/middlewares"
	"github.com/hvuhsg/goapi/request"
	"github.com/hvuhsg/goapi/responses"
)

// Private middleware used internaly
// To restrict methods to a view just use the Methods function on the view.
type methodsMiddleware struct {
	Methods []string
}

func newMethodsMiddleware(methods []string) middlewares.Middleware {
	return &methodsMiddleware{Methods: methods}
}

func (mm *methodsMiddleware) Apply(next middlewares.AppHandler) middlewares.AppHandler {
	return func(request *request.Request) responses.Response {
		for _, method := range mm.Methods {
			if method == request.HTTPRequest.Method {
				return next(request)
			}
		}

		return responses.NewErrorResponse(http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}
