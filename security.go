package goapi

import "fmt"

type SecurityProvider interface {
	GetName() string
	GetScopes() []string
	IsAuthenticated(*Request) bool
}

type APISecurity struct {
	apiKey     string
	headerName string
}

func (APISecurity) GetName() string {
	return "api-key"
}

func (APISecurity) GetScopes() []string {
	return []string{}
}

func (sec *APISecurity) IsAuthenticated(r *Request) bool {
	defer func() {
		if r := recover(); r != nil {
			panic(fmt.Sprintf("API Header '%s' not found", sec.headerName))
		}
	}()

	return r.GetString(sec.headerName) == sec.apiKey
}
