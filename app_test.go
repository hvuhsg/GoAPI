package goapi_test

import (
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/hvuhsg/goapi"
)

func TestCreateView(t *testing.T) {
	app := goapi.GoAPI("title")

	t.Run("No methods", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("expecting panic on missing methods")
			}
		}()

		app.Path("/a").Description("test view").Parameter("t").Action(
			func(request *goapi.Request) goapi.Response { return goapi.HtmlResponse{Content: "1"} },
		)
	})

	t.Run("No description", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("expecting panic on missing description")
			}
		}()

		app.Path("/b").Methods(goapi.GET).Parameter("t").Action(
			func(request *goapi.Request) goapi.Response { return goapi.HtmlResponse{Content: "1"} },
		)
	})

	t.Run("same path", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("expecting panic on same path for two views")
			}
		}()

		app.Path("/c").Methods(goapi.GET).Description("c").Parameter("t").Action(
			func(request *goapi.Request) goapi.Response { return goapi.HtmlResponse{Content: "1"} },
		)
		app.Path("/c").Methods(goapi.GET).Description("c").Parameter("t").Action(
			func(request *goapi.Request) goapi.Response { return goapi.HtmlResponse{Content: "1"} },
		)
	})

	app.Path("/").Methods(goapi.GET).Description("test view").Parameter("t").Action(
		func(request *goapi.Request) goapi.Response {
			return goapi.HtmlResponse{Content: "1"}
		},
	)
}

func TestRunApp(t *testing.T) {
	app := goapi.GoAPI("test")
	ping := app.Path("/ping")
	ping.Methods(goapi.GET)
	ping.Description("ping pong")
	ping.Parameter("age", goapi.VRequired{}, goapi.VIsInt{}, goapi.VRange{Min: 5, Max: 25})
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
	reqString := string(respBody[:])

	if reqString != `{"age":20}` {
		t.Errorf("expecting response body to be '{\"age\":20}' got '%s'", respBody)
	}
}
