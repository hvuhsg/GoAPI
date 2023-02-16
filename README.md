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

app := goapi.GoAPI()
```

Once you have created the app instance, you can add routes to it by calling the Path() method and specifying the route path. You can also set the HTTP methods that the route supports using the Methods() method. Here is an example of a route that supports GET and POST methods:

```go
var ping := app.Path("/ping")
ping.Methods(goapi.GET, goapi.POST)
ping.Description("ping pong")
ping.Parameter("age", goapi.VIsInt{}, goapi.VRange{Min: 0, Max: 120})
ping.Action(
    func (request *goapi.Request) any {
        age := request.GetInt("age")
        return map[string] interface{} {
            "message": "pong", "age": age
        }
    }
)
```

In the above example, we have created a route at `/ping` that supports GET and POST methods. We have also added a description for the route and specified that it expects an integer parameter called "age" with a minimum value of 0 and a maximum value of 120. Finally, we have specified the action that should be performed when the route is accessed. In this case, we are returning a JSON object that contains a "message" key with the value "pong" and an "age" key with the age parameter as an integer.

You can add as many routes as you like to the app instance, and each route can have its own unique set of parameters and actions.

## Validation
GoAPI comes with built-in validators that can be used to validate input data automatically. In the above example, we used the `VIsInt` validator to ensure that the "age" parameter is an integer. We also used the `VRange` validator to ensure that the "age" parameter falls within a specified range.

If the input data fails validation, the framework will automatically return an error to the client. This means you can focus on the main logic of your application, and the framework will handle the validation for you.

## API Documentation
GoAPI can automatically generate API documentation in OpenAPI format (version 3). This makes it easy to share your API with others and integrate it with other tools that support OpenAPI.

To generate the API documentation, you can simply visit the "/docs" endpoint in your web browser. This will display a user-friendly interface that allows you to view the API schema and test the API endpoints.

## Conclusion
GoAPI is a fast and easy-to-use web framework for building APIs in Go. It comes with built-in validators and can automatically generate API documentation in OpenAPI format. This means you can focus on the main logic of your application while the framework handles the validation and API documentation generation for you. GoAPI is also highly customizable, which means you can tweak its configuration to suit your needs.