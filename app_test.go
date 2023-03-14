package goapi_test

import (
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/hvuhsg/goapi"
	"github.com/hvuhsg/goapi/request"
	"github.com/hvuhsg/goapi/responses"
)

func TestCreateView(t *testing.T) {
	app := goapi.GoAPI("title", "1.0")

	t.Run("No methods", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("expecting panic on missing methods")
			}
		}()

		app.Path("/a").Description("test view").Parameter("t", goapi.QUERY).Action(
			func(request *request.Request) responses.Response { return responses.NewHTMLResponse("1", 200) },
		)
	})

	t.Run("No description", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("expecting panic on missing description")
			}
		}()

		app.Path("/b").Methods(goapi.GET).Parameter("t", goapi.QUERY).Action(
			func(request *request.Request) responses.Response { return responses.NewHTMLResponse("1", 200) },
		)
	})

	t.Run("same path", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("expecting panic on same path for two views")
			}
		}()

		app.Path("/c").Methods(goapi.GET).Description("c").Parameter("t", goapi.QUERY).Action(
			func(request *request.Request) responses.Response { return responses.NewHTMLResponse("1", 200) },
		)
		app.Path("/c").Methods(goapi.GET).Description("c").Parameter("t", goapi.QUERY).Action(
			func(request *request.Request) responses.Response { return responses.NewHTMLResponse("1", 200) },
		)
	})

	app.Path("/").Methods(goapi.GET).Description("test view").Parameter("t", goapi.QUERY).Action(
		func(request *request.Request) responses.Response {
			return responses.NewHTMLResponse("1", 200)
		},
	)
}

func TestRunApp(t *testing.T) {
	app := goapi.GoAPI("test", "1.0")
	ping := app.Path("/ping")
	ping.Methods(goapi.GET)
	ping.Description("ping pong")
	ping.Parameter("age", goapi.QUERY, goapi.VRequired{}, goapi.VIsInt{}, goapi.VRange{Min: 5, Max: 25})
	ping.Action(func(request *request.Request) responses.Response {
		return responses.NewJSONResponse(responses.Json{"age": request.GetInt("age")}, 200)
	})

	go app.Run("127.0.0.1", 8080)

	time.Sleep(time.Millisecond * 200)

	resp, err := http.Get("http://127.0.0.1:8080/ping?age=20")

	if err != nil {
		t.Error("not expecting error")
	}

	if resp.StatusCode != 200 {
		t.Errorf("expecting status-code 200 got %d", resp.StatusCode)
	}

	respBody, _ := io.ReadAll(resp.Body)
	respString := string(respBody[:])

	if respString != `{"age":20}` {
		t.Errorf("expecting response body to be '{\"age\":20}' got '%s'", respBody)
	}
}

func TestOpenAPISchema(t *testing.T) {
	app := goapi.GoAPI("example", "1.0v")
	app.Description("Example app")
	app.TermOfServiceURL("www.example.com/term_of_service")
	app.License("MIT", "")
	app.Contact("yoyo", "example.com", "goapi@example.com")
	app.Tag("math", "math operations")
	app.Tag("deprecated", "deprecated operations")

	add := app.Path("/add")
	add.Tags("math")
	add.Methods(goapi.GET)
	add.Description("Add two numbers")
	add.Parameter("a", goapi.QUERY, goapi.VRequired{}, goapi.VIsInt{}, goapi.VRange{Min: 0, Max: 100})
	add.Parameter("b", goapi.QUERY, goapi.VRequired{}, goapi.VIsInt{}, goapi.VRange{Min: 0, Max: 100})
	add.Action(func(request *request.Request) responses.Response {
		return responses.NewJSONResponse(responses.Json{"sum": request.GetInt("a") + request.GetInt("b")}, 200)
	})

	sub := app.Path("/sub")
	sub.Deprecated() // deprecated route
	sub.Tags("math", "deprecated")
	sub.Methods(goapi.GET)
	add.Description("Subtruct two numbers (a - b)")
	add.Parameter("a", goapi.QUERY, goapi.VRequired{}, goapi.VIsInt{}, goapi.VRange{Min: 0, Max: 100})
	add.Parameter("b", goapi.QUERY, goapi.VRequired{}, goapi.VIsInt{}, goapi.VRange{Min: 0, Max: 100})
	add.Action(func(request *request.Request) responses.Response {
		return responses.NewJSONResponse(responses.Json{"result": request.GetInt("a") - request.GetInt("b")}, 200)
	})

	go app.Run("127.0.0.1", 8081)

	time.Sleep(time.Millisecond * 200)

	resp, err := http.Get("http://127.0.0.1:8081/openapi.json")

	if err != nil {
		t.Error("not expecting error")
	}

	if resp.StatusCode != 200 {
		t.Errorf("expecting status-code 200 got %d", resp.StatusCode)
	}
}
