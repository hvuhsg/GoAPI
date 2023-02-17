package goapi

import (
	"fmt"
	"net/http"
)

type App struct {
	title string
	views map[string]*View
}

func GoAPI(title string) *App {
	app := new(App)
	app.title = title
	app.views = make(map[string]*View)
	return app
}

func (a *App) registerViews(mux *http.ServeMux) {
	// Register each view's path to the corresponding HTTP handler function
	for path, view := range a.views {
		mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
			if view.validMethod(r) {
				gReq := NewRequest(r)
				view.action(gReq)
				w.WriteHeader(200)
				w.Write([]byte("hello"))
			} else {
				http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			}
		})
	}
}

func (a *App) Path(path string) *View {
	_, ok := a.views[path]
	if ok {
		panic(fmt.Sprintf("path %s already registred", path))
	}

	view := NewView(path)
	a.views[path] = view

	return view
}

func (a *App) Run(host string, port int) error {
	mux := http.NewServeMux()
	a.registerViews(mux)

	addr := fmt.Sprintf("%s:%d", host, port)
	fmt.Printf("Starting server at %s\n", addr)

	return http.ListenAndServe(addr, mux)
}
