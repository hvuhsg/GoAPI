package main

import (
	"github.com/hvuhsg/goapi"
	"github.com/hvuhsg/goapi/request"
	"github.com/hvuhsg/goapi/responses"
)

func main() {
	app := goapi.GoAPI("TLS example", "1.0v")

	root := app.Path("/")
	root.Methods(goapi.GET)
	root.Description("simple route")
	root.Action(func(request *request.Request) responses.Response {
		return responses.NewHTMLResponse("<h1>HTML Over HTTPS</h1>", 200)
	})

	app.RunTLS("127.0.0.1", 8080, "./server.crt", "./server.key")
}
