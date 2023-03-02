package goapi

import (
	"fmt"
	"net/http"
)

const (
	GET     int = iota
	POST    int = iota
	PATCH   int = iota
	DELETE  int = iota
	PUT     int = iota
	HEAD    int = iota
	OPTIONS int = iota
	CONNECT int = iota
	TRACE   int = iota
)

var methodStringToCode = map[string]int{
	"GET":     GET,
	"POST":    POST,
	"PATCH":   PATCH,
	"DELETE":  DELETE,
	"PUT":     PUT,
	"HEAD":    HEAD,
	"OPTIONS": OPTIONS,
	"CONNECT": CONNECT,
	"TRACE":   TRACE,
}

var methodCodeToString = map[int]string{
	GET:     "GET",
	POST:    "POST",
	PATCH:   "PATCH",
	DELETE:  "DELETE",
	PUT:     "PUT",
	HEAD:    "HEAD",
	OPTIONS: "OPTIONS",
	CONNECT: "CONNECT",
	TRACE:   "TRACE",
}

type View struct {
	path        string
	methods     map[int]struct{}
	parameters  map[string]Parameter
	description string
	tags        []string
	depreceted  bool
	middlewares []middleware
	action      func(request *Request) Response
}

func NewView(path string) *View {
	view := new(View)
	view.path = path
	view.methods = make(map[int]struct{})
	view.parameters = make(map[string]Parameter)
	view.description = ""
	view.tags = make([]string, 0)
	view.depreceted = false
	view.middlewares = make([]middleware, 0)
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

func (v *View) validMethod(req *http.Request) bool {
	_, ok := v.methods[methodStringToCode[req.Method]]
	return ok
}

func (v *View) isValidRequest(r *Request) (bool, error) {
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

func (v *View) applyMiddlewares(appMiddlewares []middleware) {
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
		// TODO: log on DEBUG MODE
		if r := recover(); r != nil {
			fmt.Printf("ERROR: %v\n", r)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}()

	if !v.validMethod(r) {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	req := NewRequest(r)

	isValid, err := v.isValidRequest(req)

	if !isValid {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	response := v.action(req)
	w.Header().Add("Content-Type", response.contentType())
	w.WriteHeader(response.statusCode())
	w.Write(response.toBytes())
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

func (v *View) Methods(methods ...int) *View {
	for method := range methods {
		v.methods[method] = struct{}{}
	}
	return v
}

func (v *View) Description(description string) *View {
	v.requireMethods()
	v.description = description
	return v
}

func (v *View) Middlewares(middlewares ...middleware) {
	v.middlewares = append(v.middlewares, middlewares...)
}

func (v *View) Parameter(paramName string, in string, validators ...Validator) *View {
	v.requireMethods()
	v.requireDescription()
	v.parameters[paramName] = NewParameter(paramName, in, validators)
	return v
}

type AppHandler func(request *Request) Response

func (v *View) Action(r AppHandler) {
	v.requireMethods()
	v.requireDescription()

	v.action = r
}
