package goapi_test

import (
	"testing"

	"github.com/hvuhsg/goapi"
)

func TestCreateView(t *testing.T) {
	app := goapi.GoApp("title")

	t.Run("No methods", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("expecting panic on missing methods")
			}
		}()

		app.Path("/a").Description("test view").Parameter("t").Action(
			func(request any) any { return 1 },
		)
	})

	t.Run("No description", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("expecting panic on missing description")
			}
		}()

		app.Path("/b").Methods([]int{goapi.GET}).Parameter("t").Action(
			func(request any) any { return 1 },
		)
	})

	t.Run("same path", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("expecting panic on same path for two views")
			}
		}()

		app.Path("/c").Methods([]int{goapi.GET}).Description("c").Parameter("t").Action(
			func(request any) any { return 1 },
		)
		app.Path("/c").Methods([]int{goapi.GET}).Description("c").Parameter("t").Action(
			func(request any) any { return 1 },
		)
	})

	app.Path("/").Methods([]int{goapi.GET}).Description("test view").Parameter("t").Action(
		func(request any) any {
			return 1
		},
	)
}
