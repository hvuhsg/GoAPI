package responses

import "net/http"

type errorResponse struct {
	headers http.Header
	Error   string
	Code    int
}

func NewErrorResponse(error string, code int) Response {
	return errorResponse{Error: error, Code: code}
}

func (er errorResponse) Headers() http.Header {
	er.headers = http.Header{}
	er.headers.Set("Content-Type", "text/plain; charset=utf-8")
	er.headers.Set("X-Content-Type-Options", "nosniff")
	return er.headers
}

func (er errorResponse) ToBytes() []byte {
	return []byte(er.Error)
}

func (er errorResponse) StatusCode() int {
	return er.Code
}
