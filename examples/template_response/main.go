package main

import (
	"log"

	"github.com/hvuhsg/goapi"
)

func main() {
	app := goapi.GoAPI("small", "1.0v")

	echo := app.Path("/")
	echo.Methods(goapi.GET)
	echo.Description("get template")
	echo.Action(func(request *goapi.Request) goapi.Response {
		data := map[string]string{"Title": "ss", "Header": "sadd", "Content": "asdsd"}
		return goapi.TemplateResponse{TemplatePath: "./index.html", Data: data, Code: 200}
	})

	err := app.Run("127.0.0.1", 8000)
	log.Panicln(err.Error())
}
