package main

import "github.com/hvuhsg/goapi"

func main() {
	app := goapi.GoAPI("TLS example", "1.0v")

	root := app.Path("/")
	root.Methods(goapi.GET)
	root.Description("simple route")
	root.Action(func(request *goapi.Request) goapi.Response {
		return goapi.HtmlResponse{Content: "<h1>HTML Over HTTPS</h1>", Code: 200}
	})

	app.RunTLS("127.0.0.1", 8080, "./server.crt", "./server.key")
}
