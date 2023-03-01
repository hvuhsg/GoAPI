# GoAPI

<p align="center">
    <em>A Fast and Easy-to-use Web Framework for Building APIs in Go</em>
</p>
<p align="center">
	<a href="https://github.com/hvuhsg/GoAPI/actions/workflows/tests.yml" target="_blank">
		<img src="https://github.com/hvuhsg/GoAPI/actions/workflows/tests.yml/badge.svg?branch=main" alt="tests status">
	</a>
</p>

---

**Documentation**: <a href="https://github.com/hvuhsg/GoAPI/README.md" target="_blank">https://github.com/hvuhsg/GoAPI/README.md</a>

**Source Code**: <a href="https://github.com/hvuhsg/GoAPI" target="_blank">https://github.com/hvuhsg/GoAPI</a>

---


GoAPI is a web framework written in Go that is inspired by FastAPI. It allows you to quickly build and deploy RESTful APIs with minimal effort. The framework comes with built-in validators that can be used to validate and auto-build the API schema with OpenAPI format (version 3). This means you can focus on the main logic of your application while the framework handles the validation and API documentation generation.  


The key features are:
- [**automatic docs**](#api-documentation)
- [**extensible validators system**](#validation)
- [**high level syntax**](#usage)
- [**middlewares** support](#middlewares)
- [**native handlers** support](#native-handlers)

## Quick Start
### Install
To install GoAPI, you need to have Go version 1.13 or higher installed on your system. Once you have installed Go, you can use the following command to install GoAPI:

```sh
go get github.com/hvuhsg/GoAPI
```

### Usage

To use GoAPI, you need to import it in your Go code and create a new instance of the goapi.App object:

```go
package main

import "github.com/hvuhsg/goapi"

func main() {
	app := goapi.GoAPI("small", "1.0v")
}
```

Once you have created the app instance, you can add routes to it by calling the Path() method and specifying the route path. You can also set the HTTP methods that the route supports using the Methods() method. Here is an example of a route that support GET method:

```go
package main

import "github.com/hvuhsg/goapi"

func main() {
	app := goapi.GoAPI("small", "1.0v")

	echo := app.Path("/echo")
	echo.Methods(goapi.GET)
	echo.Description("echo a back")
	echo.Parameter("a", goapi.QUERY, goapi.VRequired{}, goapi.VIsInt{})
	echo.Action(func(request *goapi.Request) goapi.Response {
		return goapi.JsonResponse{"a": request.GetInt("a")}
	})

	app.Run("127.0.0.1", 8080)
}
```

In the above example, we have created a route at `/echo` that supports GET method. We have also added a description for the route and specified that it expects an integer parameter called "a". Finally, we have specified the action that should be performed when the route is accessed. In this case, we are returning a JSON object that contains a "a" key with the value of a as an integer.

You can add as many routes as you like to the app instance, and each route can have its own unique set of parameters and validators.

## Examples

To help you get started with using GoAPI, we have provided some examples in the examples directory of the repository. These examples demonstrate various use cases of the framework and how to use its features.

The examples included are:  
- **smallest** The smallest example of ready to run api.
- **math_api**: Simple use of methods, parameters and validators.  

To run the examples, navigate to the examples directory and run the following command:

```sh
go run <example_name>/main.go
```

This will start the example server `127.0.0.1:8080` and you can visit the example endpoints in your web browser or via curl.
To see all of the endpoints tou can visit `127.0.0.1:8080/docs` to see the auto generated interactive docs.

Feel free to use these examples as a starting point for your own projects and modify them as needed.

## Validation
GoAPI comes with built-in validators that can be used to validate input data automatically. In the above example, we used the `VIsInt` validator to ensure that the "timestamp" parameter is an integer. We also used the `VRange` validator to ensure that the "a" and "b" parameters falls within a specified range.

If the input data fails validation, the framework will automatically return an error to the client. This means you can focus on the main logic of your application, and the framework will handle the validation for you.

Here is a validator for example
<details>
<summary>See validator</summary>

```go
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
</details>

## Middlewares
GoAPI supports middlewares, middlewares can be defind in the app level or in the view level, middlewares in the app level are applied to all views.

here is an example for logging middleware, that logs requests in the format that nginx uses.
<details>
<summary>See middleware</summary>

```go
type SimpleLoggingMiddleware struct{}

func (SimpleLoggingMiddleware) Apply(next AppHandler) AppHandler {
	return func(request *Request) Response {
		response := next(request)

		scheme := "http"
		if request.HTTPRequest.TLS != nil {
			scheme = "https"
		}
		fullURL := fmt.Sprintf("%s://%s%s", scheme, request.HTTPRequest.Host, request.HTTPRequest.URL.String())

		method := request.HTTPRequest.Method
		path := request.HTTPRequest.URL.Path
		responseSize := len(response.toBytes())
		remoteAddr := request.HTTPRequest.RemoteAddr
		date := time.Now().Format("2006-01-02 15:04:05")
		userAgent := request.HTTPRequest.UserAgent()

		// TODO: FIX status code
		fmt.Printf("%s - - [%s] \"%s %s %s\" %d %d \"%s\" \"%s\"\n", remoteAddr, date, method, path, request.HTTPRequest.Proto, 200, responseSize, fullURL, userAgent)
		return response
	}
}

```
</details>

To apply middlewares simply pass it to the Middlewares method for the app or for the view.

```go
package main

import (
	"net/http"

	"github.com/hvuhsg/goapi"
)

func main() {
	app := goapi.GoAPI("external handler", "1.0v")
	app.Middlewares(goapi.SimpleLoggingMiddleware{})
}
```

## Native handlers
To allow the usage of native handlers we added a simple way to include them in the app, simply pass the native Handler into the Include method of the app.

Note that native handlers will not be shown in the automatic docs, and currently will not be affected by the middlewares.

here is an example:
```go
package main

import (
	"net/http"

	"github.com/hvuhsg/goapi"
)

func main() {
	app := goapi.GoAPI("external handler", "1.0v")

	// Serve files from GoAPI app using the Include method that allow us to include external handlers in the app
	app.Include("/", http.FileServer(http.Dir(".")))

	app.Run("127.0.0.1", 8080)
}
```

## API Documentation
GoAPI can automatically generate API documentation in OpenAPI format (version 3). This makes it easy to share your API with others and integrate it with other tools that support OpenAPI.

To generate the API documentation, you can simply visit the "/docs" endpoint in your web browser. This will display a user-friendly interface that allows you to view the API schema and test the API endpoints.
For the JSON schema you can visit "/openapi.json".  


![Swagger UI](/docs/images/openapi_closed.png)
![Swagger route open](/docs/images/openapi_open.png)
