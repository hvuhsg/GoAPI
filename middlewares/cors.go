package middlewares

import (
	"net/http"
	"strings"

	"github.com/hvuhsg/goapi/request"
	"github.com/hvuhsg/goapi/responses"
)

type corsMiddleware struct {
	allowedOrigins []string
	allowedMethods []string
	allowedHeaders []string
}

// Creates a new CORS middleware that allows requests from any origin (using the * character)
// and with the specified HTTP methods and headers.
func NewCORSMiddleware(allowedOrigins, allowedMethods, allowedHeaders []string) Middleware {
	return &corsMiddleware{
		allowedOrigins: allowedOrigins,
		allowedMethods: allowedMethods,
		allowedHeaders: allowedHeaders,
	}
}

func (cm *corsMiddleware) originIsAllowed(origin string) bool {
	for _, allowedOrigin := range cm.allowedOrigins {
		if allowedOrigin == "*" || allowedOrigin == origin {
			return true
		}
	}

	return false
}

func (cm corsMiddleware) Apply(next AppHandler) AppHandler {
	return func(request *request.Request) responses.Response {
		isOptionsMethod := request.HTTPRequest.Method == http.MethodOptions

		if isOptionsMethod {
			if cm.originIsAllowed(request.HTTPRequest.Header.Get("Origin")) {
				response := responses.NewHTMLResponse("", http.StatusOK)
				response.Headers().Set("Access-Control-Allow-Origin", request.HTTPRequest.Header.Get("Origin"))
				response.Headers().Set("Access-Control-Allow-Methods", strings.Join(cm.allowedMethods, ", "))
				response.Headers().Set("Access-Control-Allow-Headers", strings.Join(cm.allowedHeaders, ", "))
				return response
			} else {
				response := responses.NewErrorResponse(http.StatusText(http.StatusForbidden), http.StatusForbidden)
				return response
			}
		}

		response := next(request)
		response.Headers().Set("Access-Control-Allow-Origin", request.HTTPRequest.Header.Get("Origin"))
		response.Headers().Set("Access-Control-Allow-Methods", strings.Join(cm.allowedMethods, ", "))
		response.Headers().Set("Access-Control-Allow-Headers", strings.Join(cm.allowedHeaders, ", "))
		return response
	}
}
