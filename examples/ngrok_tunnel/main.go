package main

import (
	"log"
	"os"

	"github.com/hvuhsg/goapi"
)

func main() {
	app := goapi.GoAPI("ngrok tunnel", "1.0v")

	echo := app.Path("/")
	echo.Methods(goapi.GET)
	echo.Description("home page")
	echo.Action(func(request *goapi.Request) goapi.Response {
		return goapi.HtmlResponse{Content: "<h1>Served over ngrok tunnel</h1>"}
	})

	ngrokToken := os.Getenv("NGROK_TOKEN")
	if ngrokToken == "" {
		log.Println("Enviroment variable NGROK_TOKEN is required")
		os.Exit(1)
	}

	app.RunNgrok(ngrokToken)
}
