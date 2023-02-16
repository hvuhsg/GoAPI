package goapi

import "net/http"

// Middleware type
type Middleware func(http.Handler) http.Handler
