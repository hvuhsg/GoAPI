package middlewares

import (
	"net/http"
	"strings"

	"github.com/hvuhsg/goapi/request"
	"github.com/hvuhsg/goapi/responses"
)

type IPFilterMiddleware struct {
	AllowedIPs []string
}

func (ipm IPFilterMiddleware) Apply(next AppHandler) AppHandler {
	return func(request *request.Request) responses.Response {
		clientIP := strings.Split(request.HTTPRequest.RemoteAddr, ":")[0]

		// Check if the client IP is allowed
		if !ipm.isAllowedIP(clientIP) {
			return responses.NewHTMLResponse("Access denied", http.StatusForbidden)
		}

		return next(request)
	}
}

func (ipm IPFilterMiddleware) isAllowedIP(ip string) bool {
	for _, allowedIP := range ipm.AllowedIPs {
		if allowedIP == ip {
			return true
		}
	}
	return false
}
