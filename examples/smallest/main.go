package main

import "github.com/hvuhsg/goapi"

func main() {
	app := goapi.GoAPI("small", "1.0v")
	app.Middlewares(goapi.SimpleLoggingMiddleware{})

	echo := app.Path("/echo")
	echo.Methods(goapi.GET)
	echo.Description("echo a back")
	echo.Parameter("a", goapi.QUERY, goapi.VRequired{}, goapi.VIsInt{})
	echo.Action(func(request *goapi.Request) goapi.Response {
		return goapi.JsonResponse{"a": request.GetInt("a")}
	})

	app.Run("127.0.0.1", 8080)
}
