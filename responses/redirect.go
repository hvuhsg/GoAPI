package responses

import (
	"net/http"
	urlpkg "net/url"
	"path"
	"strconv"
	"strings"
	"unicode/utf8"
)

func hexEscapeNonASCII(s string) string {
	newLen := 0
	for i := 0; i < len(s); i++ {
		if s[i] >= utf8.RuneSelf {
			newLen += 3
		} else {
			newLen++
		}
	}
	if newLen == len(s) {
		return s
	}
	b := make([]byte, 0, newLen)
	for i := 0; i < len(s); i++ {
		if s[i] >= utf8.RuneSelf {
			b = append(b, '%')
			b = strconv.AppendInt(b, int64(s[i]), 16)
		} else {
			b = append(b, s[i])
		}
	}
	return string(b)
}

type redirectResponse struct {
	headers http.Header
	code    int
}

func newRedirectResponse(r *http.Request, url string, code int) Response {
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
			oldpath := r.URL.Path
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
	if !hadCT && (r.Method == "GET" || r.Method == "HEAD") {
		rr.headers.Set("Content-Type", "text/html; charset=utf-8")
	}

	return rr
}

func (rr redirectResponse) ToBytes() []byte {
	return nil
}

func (rr redirectResponse) StatusCode() int {
	return rr.code
}

func (rr redirectResponse) Headers() http.Header {
	return rr.headers
}

func NewPermanentRedirectResponse(r *http.Request, url string) Response {
	return newRedirectResponse(r, url, http.StatusPermanentRedirect)
}

func NewTemporaryRedirectResponse(r *http.Request, url string) Response {
	return newRedirectResponse(r, url, http.StatusTemporaryRedirect)
}
