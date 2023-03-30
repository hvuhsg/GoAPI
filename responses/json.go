package responses

import (
	"encoding/json"
	"net/http"
)

type Json map[string]any

type jsonResponse struct {
	headers http.Header
	Content Json
	Code    int
}

func NewJSONResponse(content Json, code int) Response {
	return jsonResponse{Content: content, Code: code}
}

func (jr jsonResponse) Headers() http.Header {
	jr.headers = http.Header{}
	jr.headers.Set("Content-Type", "application/json")
	return jr.headers
}

func (jr jsonResponse) ToBytes() []byte {
	bytes, err := json.Marshal(jr.Content)
	if err != nil {
		panic(err)
	}
	return bytes
}

func (jr jsonResponse) StatusCode() int {
	if jr.Code == 0 {
		jr.Code = 200
	}
	return jr.Code
}
