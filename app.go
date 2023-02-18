package goapi

import (
	"fmt"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
)

type App struct {
	title            string
	version          string
	description      string
	termOfServiceURL string
	contact          openapi3.Contact
	views            map[string]*View
}

func GoAPI(title string, version string) *App {
	app := new(App)
	app.title = title
	app.version = version
	app.description = ""
	app.termOfServiceURL = ""
	app.contact = openapi3.Contact{}
	app.views = make(map[string]*View)
	return app
}

func (a *App) openapi3Schema() ([]byte, error) {
	paths := make(openapi3.Paths)

	for _, view := range a.views {
		path, ok := paths[view.path]
		if !ok {
			path = &openapi3.PathItem{}
		}

		for method := range view.methods {
			methodString := methodCodeToString[method]
			parameters := make(openapi3.Parameters, 0)

			for paramName, paramInfo := range view.parameters {
				schemaVal := openapi3.NewSchema()
				required := false

				for _, validator := range paramInfo.validators {
					_, ok = validator.(VRequired)
					if ok {
						required = true
					}

					validator.updateOpenAPISchema(schemaVal)
				}

				paramRef := openapi3.ParameterRef{Value: &openapi3.Parameter{
					Name:       paramName,
					In:         paramInfo.in,
					Required:   required,
					Schema:     openapi3.NewSchemaRef("", schemaVal),
					Deprecated: false,
				}}
				parameters = append(parameters, &paramRef)
			}

			responses := openapi3.NewResponses()

			// Set default validation error response
			respDesc := "Error when validating request against validators"
			responses["422"] = &openapi3.ResponseRef{Value: &openapi3.Response{Description: &respDesc}}

			operation := openapi3.Operation{
				Description: view.description,
				Tags:        view.tags,
				Parameters:  parameters,
				Responses:   responses,
				Deprecated:  view.depreceted,
			}
			path.SetOperation(methodString, &operation)
		}

		paths[view.path] = path
	}

	schemaObj := openapi3.T{
		OpenAPI: "3.0.0",
		Info: &openapi3.Info{
			Title:          a.title,
			Version:        a.version,
			Description:    a.description,
			TermsOfService: a.termOfServiceURL,
			Contact:        &a.contact,
		},
		Paths: paths,
	}

	return schemaObj.MarshalJSON()
}

func (a *App) registerViews(mux *http.ServeMux) {
	// Register each view's path to the corresponding HTTP handler function
	for path, view := range a.views {
		mux.HandleFunc(path, view.requestHandler)
	}
}

func (a *App) registerInternalViews(mux *http.ServeMux) {
	// Docs internal view, retunrs the OpenAPI-3 schmea
	schema, schemaErr := a.openapi3Schema() // Only marshel on startup for performence
	mux.HandleFunc("/docs", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		if schemaErr != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.Write(schema)
	})
}

func (a *App) Description(description string) {
	a.description = description
}

func (a *App) TermOfServiceURL(termOfServiceURL string) {
	a.termOfServiceURL = termOfServiceURL
}

func (a *App) Contact(name string, url string, email string) {
	a.contact = openapi3.Contact{Name: name, URL: url, Email: email}
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
	a.registerInternalViews(mux)

	addr := fmt.Sprintf("%s:%d", host, port)
	fmt.Printf("Starting server at %s\n", addr)

	return http.ListenAndServe(addr, mux)
}
