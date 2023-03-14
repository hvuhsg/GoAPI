package main

import (
	"log"

	"github.com/hvuhsg/goapi"
	"github.com/hvuhsg/goapi/request"
	"github.com/hvuhsg/goapi/responses"
)

func main() {
	app := goapi.GoAPI("small", "1.0v")

	echo := app.Path("/")
	echo.Methods(goapi.GET)
	echo.Description("get template")
	echo.Action(func(request *request.Request) responses.Response {
		data := map[string]string{"Title": "ss", "Header": "sadd", "Content": "asdsd"}
		return responses.NewTemplateResponse("./index.html", data, 200)
	})

	err := app.Run("127.0.0.1", 8000)
	log.Panicln(err.Error())
}
