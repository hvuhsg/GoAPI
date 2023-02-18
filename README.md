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
app := goapi.GoAPI("app name")
ping := app.Path("/ping")
ping.Methods(goapi.GET)
ping.Description("Check server availability")
ping.Parameter("timestamp", goapi.VRequired{}, goapi.VIsInt{}, goapi.VRange{Min: 0, Max: 999999999999})
ping.Action(func(request *goapi.Request) goapi.Response {
    return goapi.JsonResponse{"timestamp": request.GetInt("timestamp")}
})
```

In the above example, we have created a route at `/ping` that supports GET method. We have also added a description for the route and specified that it expects an integer parameter called "timestamp" with a minimum and maximum values. Finally, we have specified the action that should be performed when the route is accessed. In this case, we are returning a JSON object that contains a "timestamp" key with the value of the timestamp parameter as an integer.

You can add as many routes as you like to the app instance, and each route can have its own unique set of parameters and validators.

## Validation
GoAPI comes with built-in validators that can be used to validate input data automatically. In the above example, we used the `VIsInt` validator to ensure that the "timestamp" parameter is an integer. We also used the `VRange` validator to ensure that the "timestamp" parameter falls within a specified range.

If the input data fails validation, the framework will automatically return an error to the client. This means you can focus on the main logic of your application, and the framework will handle the validation for you.

NOTE: the order of the validators is the order the frameword is validating them, if you want your parameter to be required and int the order must be VRequired and then VIsInt, if the order is reversed the validation will allow non int parameters to pass. (Will probebly change in later version)

Adding a custom validator is very easy, just implement the Validator interfce:
```go
type Validator interface {
	Validate(r *Request, paramName string) error
}

// VRange implementaion for refernce
type VRange struct {
	Min float64
	Max float64
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
