package goapi_test

import (
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/hvuhsg/goapi"
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
			func(request *goapi.Request) goapi.Response { return goapi.HtmlResponse{Content: "1"} },
		)
	})

	t.Run("No description", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("expecting panic on missing description")
			}
		}()

		app.Path("/b").Methods(goapi.GET).Parameter("t", goapi.QUERY).Action(
			func(request *goapi.Request) goapi.Response { return goapi.HtmlResponse{Content: "1"} },
		)
	})

	t.Run("same path", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("expecting panic on same path for two views")
			}
		}()

		app.Path("/c").Methods(goapi.GET).Description("c").Parameter("t", goapi.QUERY).Action(
			func(request *goapi.Request) goapi.Response { return goapi.HtmlResponse{Content: "1"} },
		)
		app.Path("/c").Methods(goapi.GET).Description("c").Parameter("t", goapi.QUERY).Action(
			func(request *goapi.Request) goapi.Response { return goapi.HtmlResponse{Content: "1"} },
		)
	})

	app.Path("/").Methods(goapi.GET).Description("test view").Parameter("t", goapi.QUERY).Action(
		func(request *goapi.Request) goapi.Response {
			return goapi.HtmlResponse{Content: "1"}
		},
	)
}

func TestRunApp(t *testing.T) {
	app := goapi.GoAPI("test", "1.0")
	ping := app.Path("/ping")
	ping.Methods(goapi.GET)
	ping.Description("ping pong")
	ping.Parameter("age", goapi.QUERY, goapi.VRequired{}, goapi.VIsInt{}, goapi.VRange{Min: 5, Max: 25})
	ping.Action(func(request *goapi.Request) goapi.Response {
		return goapi.JsonResponse{"age": request.GetInt("age")}
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
	app := goapi.GoAPI("test", "1.0")
	app.Description("App for testing")
	app.TermOfServiceURL("www.example.com/term_of_service")
	app.Contact("yoyo", "example.com", "goapi@example.com")

	ping := app.Path("/ping")
	ping.Methods(goapi.GET)
	ping.Description("ping pong")
	ping.Parameter("age", goapi.QUERY, goapi.VRequired{}, goapi.VIsInt{}, goapi.VRange{Min: 5, Max: 25})
	ping.Action(func(request *goapi.Request) goapi.Response {
		return goapi.JsonResponse{"age": request.GetInt("age")}
	})

	pong := app.Path("/pong/{age}/{password}")
	pong.Deprecated()
	pong.Tags("pong", "bla")
	pong.Methods(goapi.PATCH)
	pong.Description("bla bla")
	pong.Parameter("age", goapi.PATH, goapi.VRequired{}, goapi.VIsInt{}, goapi.VRange{Min: 5, Max: 25})
	pong.Parameter("password", goapi.PATH, goapi.VRequired{}, goapi.VStringLength{Min: 3, Max: 12})
	pong.Action(func(request *goapi.Request) goapi.Response {
		return goapi.JsonResponse{"age": request.GetInt("age"), "password": request.GetString("password")}
	})

	go app.Run("127.0.0.1", 8081)

	time.Sleep(time.Millisecond * 200)

	resp, err := http.Get("http://127.0.0.1:8081/docs")

	if err != nil {
		t.Error("not expecting error")
	}

	if resp.StatusCode != 200 {
		t.Errorf("expecting status-code 200 got %d", resp.StatusCode)
	}

	respBody, _ := io.ReadAll(resp.Body)
	respString := string(respBody[:])

	fmt.Println(respString)
}
