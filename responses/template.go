package responses

import (
	"bytes"
	"net/http"
	"text/template"
)

type templateResponse struct {
	headers      http.Header
	TemplatePath string
	Data         any
	Code         int
}

func NewTemplateResponse(tmpPath string, data any, code int) Response {
	return templateResponse{TemplatePath: tmpPath, Data: data, Code: code}
}

func (tr templateResponse) Headers() http.Header {
	tr.headers = http.Header{}
	tr.headers.Set("Content-Type", "text/html")
	return tr.headers
}

func (tr templateResponse) ToBytes() []byte {
	tpl, err := template.ParseFiles(tr.TemplatePath)
	if err != nil {
		panic("Can't parse template: " + err.Error())
	}

	var buffer bytes.Buffer
	err = tpl.Execute(&buffer, tr.Data)
	if err != nil {
		panic("Can't populate template with data: " + err.Error())
	}

	return buffer.Bytes()
}

func (tr templateResponse) StatusCode() int {
	if tr.Code == 0 {
		tr.Code = 200
	}
	return tr.Code
}
