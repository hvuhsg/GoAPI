package responses

import "net/http"

type htmlResponse struct {
	headers http.Header
	Content string
	Code    int
}

func NewHTMLResponse(content string, code int) Response {
	return htmlResponse{Content: content, Code: code}
}

func (hr htmlResponse) Headers() http.Header {
	hr.headers = http.Header{}
	hr.headers.Set("Content-Type", "text/html")
	return hr.headers
}

func (hr htmlResponse) ToBytes() []byte {
	return []byte(hr.Content)
}

func (hr htmlResponse) StatusCode() int {
	if hr.Code == 0 {
		hr.Code = 200
	}
	return hr.Code
}
