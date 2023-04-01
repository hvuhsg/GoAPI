package responses

import (
	"net/http"
)

type Response interface {
	Headers() http.Header
	ToBytes() []byte
	StatusCode() int
}

type response struct {
	Header  http.Header
	content []byte
	code    int
}

func (r *response) Headers() http.Header {
	return r.Header
}

func (r *response) ToBytes() []byte {
	return r.content
}

func (r *response) StatusCode() int {
	return r.code
}

func NewResponse(content []byte, code int) Response {
	return &response{content: content, code: code}
}
