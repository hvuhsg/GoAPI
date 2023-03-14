package middlewares

import (
	"github.com/hvuhsg/goapi/request"
	"github.com/hvuhsg/goapi/responses"
)

type Middleware interface {
	Apply(AppHandler) AppHandler
}

type AppHandler func(request *request.Request) responses.Response
