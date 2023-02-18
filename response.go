package goapi

import "encoding/json"

type Response interface {
	toBytes() []byte
	contentType() string
}

type HtmlResponse struct {
	Content string
}

func (HtmlResponse) contentType() string {
	return "text/html"
}

func (hr HtmlResponse) toBytes() []byte {
	return []byte(hr.Content)
}

type JsonResponse map[string]any

func (JsonResponse) contentType() string {
	return "application/json"
}

func (jr JsonResponse) toBytes() []byte {
	bytes, _ := json.Marshal(jr)
	return bytes
}
