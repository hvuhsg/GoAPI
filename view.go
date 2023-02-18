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

type View struct {
	path        string
	methods     map[int]struct{}
	parameters  map[string]Parameter
	description string
	action      func(request *Request) any
}

func NewView(path string) *View {
	view := new(View)
	view.path = path
	view.methods = make(map[int]struct{})
	view.parameters = make(map[string]Parameter)
	view.description = ""
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

func (v *View) requestHandler(w http.ResponseWriter, r *http.Request) {
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
	w.WriteHeader(200)
	w.Write([]byte(fmt.Sprint(response)))
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

func (v *View) Parameter(paramName string, validators ...Validator) *View {
	v.requireMethods()
	v.requireDescription()
	v.parameters[paramName] = NewParameter(paramName, validators)
	return v
}

func (v *View) Action(f func(request *Request) any) {
	v.requireMethods()
	v.requireDescription()

	v.action = f
}
