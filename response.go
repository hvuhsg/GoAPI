package goapi

import (
	"bytes"
	"encoding/json"
	"net/http"
	urlpkg "net/url"
	"path"
	"strings"
	"text/template"
)

type Response interface {
	Headers() http.Header
	toBytes() []byte
	statusCode() int
}

type HtmlResponse struct {
	headers http.Header
	Content string
	Code    int
}

func (hr HtmlResponse) Headers() http.Header {
	hr.headers = http.Header{}
	hr.headers.Add("Content-Type", "text/html")
	return hr.headers
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
	headers http.Header
	Content Json
	Code    int
}

func (jr JsonResponse) Headers() http.Header {
	jr.headers = http.Header{}
	jr.headers.Add("Content-Type", "application/json")
	return jr.headers
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

type redirectResponse struct {
	headers http.Header
	code    int
}

func newRedirectResponse(r *Request, url string, code int) Response {
	rr := new(redirectResponse)
	rr.code = code
	rr.headers = http.Header{}

	if u, err := urlpkg.Parse(url); err == nil {
		// If url was relative, make its path absolute by
		// combining with request path.
		// The client would probably do this for us,
		// but doing it ourselves is more reliable.
		// See RFC 7231, section 7.1.2
		if u.Scheme == "" && u.Host == "" {
			oldpath := r.HTTPRequest.URL.Path
			if oldpath == "" { // should not happen, but avoid a crash if it does
				oldpath = "/"
			}

			// no leading http://server
			if url == "" || url[0] != '/' {
				// make relative path absolute
				olddir, _ := path.Split(oldpath)
				url = olddir + url
			}

			var query string
			if i := strings.Index(url, "?"); i != -1 {
				url, query = url[:i], url[i:]
			}

			// clean up but preserve trailing slash
			trailing := strings.HasSuffix(url, "/")
			url = path.Clean(url)
			if trailing && !strings.HasSuffix(url, "/") {
				url += "/"
			}
			url += query
		}
	}

	// RFC 7231 notes that a short HTML body is usually included in
	// the response because older user agents may not understand 301/307.
	// Do it only if the request didn't already have a Content-Type header.
	_, hadCT := rr.headers["Content-Type"]

	rr.headers.Set("Location", hexEscapeNonASCII(url))
	if !hadCT && (r.HTTPRequest.Method == "GET" || r.HTTPRequest.Method == "HEAD") {
		rr.headers.Set("Content-Type", "text/html; charset=utf-8")
	}

	return rr
}

func (rr redirectResponse) toBytes() []byte {
	return nil
}

func (rr redirectResponse) statusCode() int {
	return rr.code
}

func (rr redirectResponse) Headers() http.Header {
	return rr.headers
}

func NewPermanentRedirectResponse(r *Request, url string) Response {
	return newRedirectResponse(r, url, http.StatusPermanentRedirect)
}

func NewTemporaryRedirectResponse(r *Request, url string) Response {
	return newRedirectResponse(r, url, http.StatusTemporaryRedirect)
}

type TemplateResponse struct {
	headers      http.Header
	TemplatePath string
	Data         interface{}
	Code         int
}

func (tr TemplateResponse) Headers() http.Header {
	tr.headers = http.Header{}
	tr.headers.Add("Content-Type", "text/html")
	return tr.headers
}

func (tr TemplateResponse) toBytes() []byte {
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

func (tr TemplateResponse) statusCode() int {
	if tr.Code == 0 {
		tr.Code = 200
	}
	return tr.Code
}
