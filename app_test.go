package goapi_test

import (
	"fmt"
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
			func(request *goapi.Request) any { return 1 },
		)
	})

	t.Run("No description", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("expecting panic on missing description")
			}
		}()

		app.Path("/b").Methods(goapi.GET).Parameter("t").Action(
			func(request *goapi.Request) any { return 1 },
		)
	})

	t.Run("same path", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("expecting panic on same path for two views")
			}
		}()

		app.Path("/c").Methods(goapi.GET).Description("c").Parameter("t").Action(
			func(request *goapi.Request) any { return 1 },
		)
		app.Path("/c").Methods(goapi.GET).Description("c").Parameter("t").Action(
			func(request *goapi.Request) any { return 1 },
		)
	})

	app.Path("/").Methods(goapi.GET).Description("test view").Parameter("t").Action(
		func(request *goapi.Request) any {
			return 1
		},
	)
}

func TestRunApp(t *testing.T) {
	app := goapi.GoAPI("test")
	ping := app.Path("/ping")
	ping.Methods(goapi.GET)
	ping.Description("ping pong")
	ping.Parameter("age", goapi.VIsInt{}, goapi.VRange{Min: 5, Max: 25})
	ping.Action(func(request *goapi.Request) any {
		fmt.Println(request.GetInt("age"))
		return 1
	})

	go app.Run("127.0.0.1", 8080)

	time.Sleep(time.Millisecond * 200)

	resp, _ := http.Get("http://127.0.0.1:8080/ping?age=2")
	fmt.Printf("response: %d", resp.StatusCode)
}
