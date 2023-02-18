# GoAPI - A Fast and Easy-to-use Web Framework for Building APIs in Go
GoAPI is a web framework written in Go that is inspired by FastAPI. It allows you to quickly build and deploy RESTful APIs with minimal effort. The framework comes with built-in validators that can be used to validate and auto-build the API schema with OpenAPI format (version 3). This means you can focus on the main logic of your application while the framework handles the validation and API documentation generation.


## Quick Start
### Install
To install GoAPI, you need to have Go version 1.13 or higher installed on your system. Once you have installed Go, you can use the following command to install GoAPI:

```sh
go get github.com/hvuhsg/GoAPI
```

### Usage

To use GoAPI, you need to import it in your Go code and create a new instance of the goapi.App object:

```go
import "github.com/hvuhsg/GoAPI"

app := goapi.GoAPI("example", "1.0v")
app.Description("Example app")                           // optional
app.TermOfServiceURL("www.example.com/term_of_service")  // optional
app.Contact("yoyo", "example.com", "goapi@example.com")  // optional
```

Once you have created the app instance, you can add routes to it by calling the Path() method and specifying the route path. You can also set the HTTP methods that the route supports using the Methods() method. Here is an example of a routes that support GET method:

```go
add := app.Path("/add")
add.Tags("math")  // optional
add.Methods(goapi.GET)
add.Description("Add two numbers")
add.Parameter("a", goapi.QUERY, goapi.VRequired{}, goapi.VIsInt{}, goapi.VRange{Min: 0, Max: 100})
add.Parameter("b", goapi.QUERY, goapi.VRequired{}, goapi.VIsInt{}, goapi.VRange{Min: 0, Max: 100})
add.Action(func(request *goapi.Request) goapi.Response {
    return goapi.JsonResponse{"sum": request.GetInt("a") + request.GetInt("b")}
})

sub := app.Path("/sub")
sub.Deprecated()                // mark as deprecated route
sub.Tags("math", "deprecated")  // optional
sub.Methods(goapi.GET)
add.Description("Subtruct two numbers (a - b)")
add.Parameter("a", goapi.QUERY, goapi.VRequired{}, goapi.VIsInt{}, goapi.VRange{Min: 0, Max: 100})
add.Parameter("b", goapi.QUERY, goapi.VRequired{}, goapi.VIsInt{}, goapi.VRange{Min: 0, Max: 100})
add.Action(func(request *goapi.Request) goapi.Response {
    return goapi.JsonResponse{"result": request.GetInt("a") - request.GetInt("b")}
})
```

In the above example, we have created a route at `/add` that supports GET method. We have also added a description for the route and specified that it expects an integer parameters called "a" and "b" with a minimum and maximum values. Finally, we have specified the action that should be performed when the route is accessed. In this case, we are returning a JSON object that contains a "sum" key with the value of a + b as an integer.

You can add as many routes as you like to the app instance, and each route can have its own unique set of parameters and validators (just like the deprecated sub route).

## Examples

To help you get started with using GoAPI, we have provided some examples in the examples directory of the repository. These examples demonstrate various use cases of the framework and how to use its features.

To run the examples, navigate to the examples directory and run the following command:

```sh
go run <example_name>/main.go
```

This will start the example server `127.0.0.1:8080` and you can visit the example endpoints in your web browser or via curl.
To see all of the endpoints tou can visit `127.0.0.1:8080/docs` to see the auto generated interactive docs.

The examples included are:  

- **math_api**: Simple use of methods, parameters and validators.  

Feel free to use these examples as a starting point for your own projects and modify them as needed.

## Validation
GoAPI comes with built-in validators that can be used to validate input data automatically. In the above example, we used the `VIsInt` validator to ensure that the "timestamp" parameter is an integer. We also used the `VRange` validator to ensure that the "a" and "b" parameters falls within a specified range.

If the input data fails validation, the framework will automatically return an error to the client. This means you can focus on the main logic of your application, and the framework will handle the validation for you.

Adding a custom validator is very easy, just implement the Validator interfce:
```go
type Validator interface {
	Validate(r *Request, paramName string) error
	updateOpenAPISchema(schema *openapi3.Schema)
}

// VRange implementaion for refernce
type VRange struct {
	Min float64
	Max float64
}

func (v VRange) updateOpenAPISchema(schema *openapi3.Schema) {
	schema.Min = &v.Min
	schema.Max = &v.Max
}
func (v VRange) Validate(r *Request, paramName string) error {
	vr := VIsFloat{}
	err := vr.Validate(r, paramName)
	if err != nil {
		return err
	}

	fVal := r.GetFloat(paramName)

	if fVal < v.Min || fVal > v.Max {
		return fmt.Errorf("parameter %s must be between %f and %f", paramName, v.Min, v.Max)
	}

	return nil
}
```

## API Documentation
GoAPI can automatically generate API documentation in OpenAPI format (version 3). This makes it easy to share your API with others and integrate it with other tools that support OpenAPI.

To generate the API documentation, you can simply visit the "/docs" endpoint in your web browser. This will display a user-friendly interface that allows you to view the API schema and test the API endpoints.
For the JSON schema you can visit "/openapi.json".  


![Swagger UI](/docs/images/openapi_closed.png)
![Swagger route open](/docs/images/openapi_open.png)
