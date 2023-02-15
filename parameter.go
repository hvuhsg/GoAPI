package goapi

type Parameter struct {
	name       string
	validators []Validator
}

func NewParameter(name string, validators []Validator) Parameter {
	return Parameter{name: name, validators: validators}
}
