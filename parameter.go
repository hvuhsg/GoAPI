package goapi

const (
	PATH   = "path"
	QUERY  = "query"
	HEADER = "header"
	COOKIE = "COOKIE"
)

type Parameter struct {
	name       string
	in         string // Where can we find this parameter (QUERY, PATH, HEADER, COOKIE)
	validators []Validator
}

func NewParameter(name string, in string, validators []Validator) Parameter {
	return Parameter{name: name, in: in, validators: validators}
}
