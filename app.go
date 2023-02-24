package goapi

import (
	"fmt"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
)

// App represents the main application.
type App struct {
	title            string
	version          string
	description      string
	termOfServiceURL string
	contact          openapi3.Contact
	tags             openapi3.Tags
	views            map[string]*View // A map of View objects keyed by their URL paths
	openapiDocsURL   string           // URL path for the OpenAPI documentation
	openapiSchemaURL string           // URL path for the OpenAPI schema
}

// GoAPI creates a new instance of the App.
func GoAPI(title string, version string) *App {
	app := new(App)
	app.title = title
	app.version = version
	app.description = ""
	app.termOfServiceURL = ""
	app.contact = openapi3.Contact{}
	app.tags = openapi3.Tags{}
	app.views = make(map[string]*View)
	app.openapiDocsURL = "/docs"
	app.openapiSchemaURL = "/openapi.json"
	return app
}

// registerViews registers each View's path to its corresponding HTTP handler function.
func (a *App) registerViews(mux *http.ServeMux) {
	for path, view := range a.views {
		mux.HandleFunc(path, view.requestHandler)
	}
}

// registerInternalViews registers internal views, such as the OpenAPI documentation route.
func (a *App) registerInternalViews(mux *http.ServeMux) {
	registerDocs(a, mux) // register OpenAPI documentation route
}

// Description sets the application description.
func (a *App) Description(description string) {
	a.description = description
}

// TermOfServiceURL sets the URL for the application's terms of service.
func (a *App) TermOfServiceURL(termOfServiceURL string) {
	a.termOfServiceURL = termOfServiceURL
}

// Contact sets the application's contact information.
func (a *App) Contact(name string, url string, email string) {
	a.contact = openapi3.Contact{Name: name, URL: url, Email: email}
}

func (a *App) Tag(name string, description string) {
	a.tags = append(a.tags, &openapi3.Tag{Name: name, Description: description})
}

// OpenapiDocsURL sets the URL path for the OpenAPI documentation.
func (a *App) OpenapiDocsURL(docsUrl string) {
	a.openapiDocsURL = docsUrl
}

// OpenapiSchemaURL sets the URL path for the OpenAPI schema.
func (a *App) OpenapiSchemaURL(schemaUrl string) {
	a.openapiSchemaURL = schemaUrl
}

// Path creates a new View for the given URL path and adds it to the App.
func (a *App) Path(path string) *View {
	_, ok := a.views[path]
	if ok {
		panic(fmt.Sprintf("path %s already registered", path))
	}

	view := NewView(path)
	a.views[path] = view

	return view
}

// Run starts the application and listens for incoming requests.
func (a *App) Run(host string, port int) error {
	mux := http.NewServeMux()
	a.registerViews(mux)
	a.registerInternalViews(mux)

	addr := fmt.Sprintf("%s:%d", host, port)
	fmt.Printf("Starting server at %s\n", addr)

	return http.ListenAndServe(addr, mux)
}
