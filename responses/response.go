package responses

import (
	"net/http"
)

type Response interface {
	Headers() http.Header
	ToBytes() []byte
	StatusCode() int
}
