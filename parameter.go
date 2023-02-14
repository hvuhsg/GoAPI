package goapi

type Parameter struct {
	name string
	// TODO: add validator
}

func NewParameter(name string) *Parameter {
	param := new(Parameter)
	param.name = name
	return param
}
