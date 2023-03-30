package goapi

import (
	"log"
	"net/http"

	"github.com/hvuhsg/goapi/middlewares"
	"github.com/hvuhsg/goapi/request"
	"github.com/hvuhsg/goapi/responses"
	"github.com/hvuhsg/goapi/validators"
)

const (
	GET     string = http.MethodGet
	POST    string = http.MethodPost
	PATCH   string = http.MethodPatch
	DELETE  string = http.MethodDelete
	PUT     string = http.MethodPut
	HEAD    string = http.MethodHead
	OPTIONS string = http.MethodOptions
	CONNECT string = http.MethodConnect
	TRACE   string = http.MethodTrace
)

type View struct {
	path        string
	methods     []string
	parameters  map[string]Parameter
	description string
	tags        []string
	depreceted  bool
	middlewares []middlewares.Middleware
	action      func(request *request.Request) responses.Response
}

func NewView(path string) *View {
	view := new(View)
	view.path = path
	view.methods = make([]string, 0)
	view.parameters = make(map[string]Parameter)
	view.description = ""
	view.tags = make([]string, 0)
	view.depreceted = false
	view.middlewares = make([]middlewares.Middleware, 0)
	view.action = nil
	return view
}

func (v *View) requireMethods() {
	if len(v.methods) == 0 {
		panic("declere methods first '.Methods[GET, POST, ...]'")
	}
}

func (v *View) requireDescription() {
	if v.description == "" {
		panic("declere description for view '.Description(\"<view description>\")'")
	}
}

func (v *View) isValidRequest(r *request.Request) (bool, error) {
	for paramName, param := range v.parameters {
		for _, validator := range param.validators {
			err := validator.Validate(r, paramName)
			if err != nil {
				return false, err
			}
		}
	}

	return true, nil
}

func (v *View) applyMiddlewares(appMiddlewares []middlewares.Middleware) {
	// Add methods middleware
	mm := newMethodsMiddleware(v.methods)
	v.action = mm.Apply(v.action)

	// Apply app middlewares
	for i := len(appMiddlewares) - 1; i >= 0; i-- {
		m := appMiddlewares[i]
		v.action = m.Apply(v.action)
	}

	// Apply view middlewares
	for i := len(v.middlewares) - 1; i >= 0; i-- {
		m := v.middlewares[i]
		v.action = m.Apply(v.action)
	}
}

func (v *View) requestHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		// If paniced; responde with 500 internal server error
		if r := recover(); r != nil {
			log.Printf("ERROR: %v\n", r)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}()

	req := request.NewRequest(r)

	isValid, err := v.isValidRequest(req)
	if !isValid {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	response := v.action(req)

	// copy response headers to response writer
	for k, values := range response.Headers() {
		for _, value := range values {
			w.Header().Add(k, value)
		}
	}

	w.WriteHeader(response.StatusCode())
	w.Write(response.ToBytes())
}

// Mark view as deprecated
func (v *View) Deprecated() *View {
	v.depreceted = true
	return v
}

func (v *View) Tags(tags ...string) *View {
	v.tags = tags
	return v
}

func (v *View) Methods(methods ...string) *View {
	v.methods = methods
	return v
}

func (v *View) Description(description string) *View {
	v.requireMethods()
	v.description = description
	return v
}

func (v *View) Middlewares(middlewares ...middlewares.Middleware) {
	v.middlewares = append(v.middlewares, middlewares...)
}

func (v *View) Parameter(paramName string, in string, validators ...validators.Validator) *View {
	v.requireMethods()
	v.requireDescription()
	v.parameters[paramName] = NewParameter(paramName, in, validators)
	return v
}

type AppHandler func(request *request.Request) responses.Response

func (v *View) Action(r AppHandler) {
	v.requireMethods()
	v.requireDescription()

	v.action = r
}
