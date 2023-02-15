package goapi_test

import (
	"testing"

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
