package main

import (
	"github.com/hvuhsg/goapi"
	"github.com/hvuhsg/goapi/request"
	"github.com/hvuhsg/goapi/responses"
	"github.com/hvuhsg/goapi/validators"
)

func main() {
	app := goapi.GoAPI("redirect", "1.0v")

	echo := app.Path("/")
	echo.Methods(goapi.GET)
	echo.Description("redirect to path")
	echo.Parameter("to", goapi.QUERY, validators.VRequired{}, validators.VIsString{})
	echo.Action(func(request *request.Request) responses.Response {
		return responses.NewPermanentRedirectResponse(request.HTTPRequest, "https://"+request.GetString("to"))
	})

	app.Run("127.0.0.1", 8080)
}
