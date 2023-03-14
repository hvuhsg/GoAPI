package main

import (
	"github.com/hvuhsg/goapi"
	"github.com/hvuhsg/goapi/request"
	"github.com/hvuhsg/goapi/responses"
	"github.com/hvuhsg/goapi/validators"
)

func main() {
	app := goapi.GoAPI("small", "1.0v")

	echo := app.Path("/echo")
	echo.Methods(goapi.GET)
	echo.Description("echo a back")
	echo.Parameter("a", goapi.QUERY, validators.VRequired{}, validators.VIsInt{})
	echo.Action(func(request *request.Request) responses.Response {
		return responses.NewJSONResponse(responses.Json{"a": request.GetInt("a")}, 200)
	})

	app.Run("127.0.0.1", 8080)
}
