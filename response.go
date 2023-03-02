package goapi

import "encoding/json"

type Response interface {
	toBytes() []byte
	contentType() string
	statusCode() int
}

type HtmlResponse struct {
	Content string
	Code    int
}

func (HtmlResponse) contentType() string {
	return "text/html"
}

func (hr HtmlResponse) toBytes() []byte {
	return []byte(hr.Content)
}

func (hr HtmlResponse) statusCode() int {
	if hr.Code == 0 {
		hr.Code = 200
	}
	return hr.Code
}

type Json map[string]any

type JsonResponse struct {
	Content Json
	Code    int
}

func (JsonResponse) contentType() string {
	return "application/json"
}

func (jr JsonResponse) toBytes() []byte {
	bytes, err := json.Marshal(jr.Content)
	if err != nil {
		panic(err)
	}
	return bytes
}

func (jr JsonResponse) statusCode() int {
	if jr.Code == 0 {
		jr.Code = 200
	}
	return jr.Code
}
