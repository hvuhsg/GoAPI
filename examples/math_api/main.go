package main

import (
	"fmt"

	"github.com/hvuhsg/goapi"
	"github.com/hvuhsg/goapi/request"
	"github.com/hvuhsg/goapi/responses"
	"github.com/hvuhsg/goapi/validators"
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
	add.Parameter("a", goapi.QUERY, validators.VRequired{}, validators.VIsInt{}, validators.VRange{Min: 0, Max: 100})
	add.Parameter("b", goapi.QUERY, validators.VRequired{}, validators.VIsInt{}, validators.VRange{Min: 0, Max: 100})
	add.Action(func(request *request.Request) responses.Response {
		return responses.NewJSONResponse(responses.Json{"sum": request.GetInt("a") + request.GetInt("b")}, 200)
	})

	sub := app.Path("/sub")
	sub.Tags("math")
	sub.Methods(goapi.GET)
	sub.Description("Subtruct two numbers (a - b)")
	sub.Parameter("a", goapi.QUERY, validators.VRequired{}, validators.VIsInt{}, validators.VRange{Min: 0, Max: 100})
	sub.Parameter("b", goapi.QUERY, validators.VRequired{}, validators.VIsInt{}, validators.VRange{Min: 0, Max: 100})
	sub.Action(func(request *request.Request) responses.Response {
		return responses.NewJSONResponse(responses.Json{"result": request.GetInt("a") - request.GetInt("b")}, 200)
	})

	mul := app.Path("/mul")
	mul.Tags("math", "array")
	mul.Methods(goapi.PUT, goapi.CONNECT)
	mul.Description("Multiply array of numbers")
	mul.Parameter("numbers", goapi.QUERY, validators.VRequired{}, validators.VIsArray{})
	mul.Action(func(request *request.Request) responses.Response {
		res := 1
		numbers := request.GetIntArray("numbers")
		for _, num := range numbers {
			res = res * num
		}

		return responses.NewJSONResponse(responses.Json{"result": res}, 200)
	})

	if err := app.Run("127.0.0.1", 8080); err != nil {
		fmt.Printf("Exiting with error: %s\n", err.Error())
	}
}
