package goapi

import "github.com/hvuhsg/goapi/validators"

const (
	PATH   = "path"
	QUERY  = "query"
	HEADER = "header"
	COOKIE = "COOKIE"
)

type Parameter struct {
	name       string
	in         string // Where can we find this parameter (QUERY, PATH, HEADER, COOKIE)
	validators []validators.Validator
}

func NewParameter(name string, in string, validators []validators.Validator) Parameter {
	return Parameter{name: name, in: in, validators: validators}
}
