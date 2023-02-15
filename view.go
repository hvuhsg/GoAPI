package goapi

import "net/http"

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
	_path        string
	_methods     map[int]struct{}
	_parameters  map[string]Parameter
	_description string
	_action      func(request *Request) any
}

func NewView(path string) *View {
	view := new(View)
	view._path = path
	view._methods = make(map[int]struct{})
	view._parameters = make(map[string]Parameter)
	view._description = ""
	view._action = nil
	return view
}

func (v *View) requireMethods() {
	if len(v._methods) == 0 {
		panic("declere methods first '.Methods[GET, POST, ...]'")
	}
}

func (v *View) requireDescription() {
	if v._description == "" {
		panic("declere description for view '.Description(\"<view description>\")'")
	}
}

func (v *View) validMethod(req *http.Request) bool {
	_, ok := v._methods[methodStringToCode[req.Method]]
	return ok
}

func (v *View) Methods(methods ...int) *View {
	for method := range methods {
		v._methods[method] = struct{}{}
	}
	return v
}

func (v *View) Description(description string) *View {
	v.requireMethods()
	v._description = description
	return v
}

func (v *View) Parameter(paramName string, validators ...Validator) *View {
	v.requireMethods()
	v.requireDescription()
	v._parameters[paramName] = NewParameter(paramName, validators)
	return v
}

func (v *View) Action(f func(request *Request) any) {
	v.requireMethods()
	v.requireDescription()

	v._action = f
}
