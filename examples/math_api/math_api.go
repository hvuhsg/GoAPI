package main

import (
	"fmt"

	"github.com/hvuhsg/goapi"
)

func main() {
	app := goapi.GoAPI("math api", "1.0v")
	app.Description("Math operation with easy api")
	app.Contact("yoyo", "mathapi.goapi.com", "mathapi@goapi.com")
	app.TermOfServiceURL("mathapi.goapi.com/terms-of-service.html")

	add := app.Path("/add")
	add.Tags("math")
	add.Methods(goapi.GET)
	add.Description("Add two numbers")
	add.Parameter("a", goapi.QUERY, goapi.VRequired{}, goapi.VIsInt{}, goapi.VRange{Min: 0, Max: 100})
	add.Parameter("b", goapi.QUERY, goapi.VRequired{}, goapi.VIsInt{}, goapi.VRange{Min: 0, Max: 100})
	add.Action(func(request *goapi.Request) goapi.Response {
		return goapi.JsonResponse{"sum": request.GetInt("a") + request.GetInt("b")}
	})

	sub := app.Path("/sub")
	sub.Tags("math")
	sub.Methods(goapi.GET)
	sub.Description("Subtruct two numbers (a - b)")
	sub.Parameter("a", goapi.QUERY, goapi.VRequired{}, goapi.VIsInt{}, goapi.VRange{Min: 0, Max: 100})
	sub.Parameter("b", goapi.QUERY, goapi.VRequired{}, goapi.VIsInt{}, goapi.VRange{Min: 0, Max: 100})
	sub.Action(func(request *goapi.Request) goapi.Response {
		return goapi.JsonResponse{"result": request.GetInt("a") - request.GetInt("b")}
	})

	mul := app.Path("/mul")
	mul.Tags("math", "array")
	mul.Methods(goapi.PUT, goapi.CONNECT)
	mul.Description("Multiply array of numbers")
	mul.Parameter("numbers", goapi.QUERY, goapi.VRequired{}, goapi.VIsArray{})
	mul.Action(func(request *goapi.Request) goapi.Response {
		res := 1
		numbers := request.GetIntArray("numbers")
		for _, num := range numbers {
			res = res * num
		}

		return goapi.JsonResponse{"result": res}
	})

	if err := app.Run("127.0.0.1", 8080); err != nil {
		fmt.Printf("Exiting with error: %s\n", err.Error())
	}
}
