package goapi

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

type middleware interface {
	Apply(AppHandler) AppHandler
}

type SimpleLoggingMiddleware struct{}

func (SimpleLoggingMiddleware) Apply(next AppHandler) AppHandler {
	return func(request *Request) Response {
		response := next(request)

		scheme := "http"
		if request.HTTPRequest.TLS != nil {
			scheme = "https"
		}
		fullURL := fmt.Sprintf("%s://%s%s", scheme, request.HTTPRequest.Host, request.HTTPRequest.URL.String())

		method := request.HTTPRequest.Method
		path := request.HTTPRequest.URL.Path
		responseSize := len(response.toBytes())
		remoteAddr := request.HTTPRequest.RemoteAddr
		date := time.Now().Format("2006-01-02 15:04:05")
		userAgent := request.HTTPRequest.UserAgent()
		statusCode := response.statusCode()

		fmt.Printf("%s - - [%s] \"%s %s %s\" %d %d \"%s\" \"%s\"\n", remoteAddr, date, method, path, request.HTTPRequest.Proto, statusCode, responseSize, fullURL, userAgent)
		return response
	}
}

type TimeoutMiddleware struct {
	Timeout time.Duration
}

func (tm TimeoutMiddleware) Apply(next AppHandler) AppHandler {
	return func(request *Request) Response {
		// Create a context with a timeout
		ctx, cancel := context.WithTimeout(context.Background(), tm.Timeout)
		defer cancel()

		// Add context to the request
		request.HTTPRequest = request.HTTPRequest.WithContext(ctx)

		// Use a goroutine to execute the request and wait for it to complete or timeout
		ch := make(chan Response, 1)
		go func() {
			ch <- next(request)
		}()

		select {
		case response := <-ch:
			return response
		case <-ctx.Done():
			return HtmlResponse{Content: "Request timed out", Code: http.StatusGatewayTimeout}
		}
	}
}

type TimingMiddleware struct{}

func (TimingMiddleware) Apply(next AppHandler) AppHandler {
	return func(request *Request) Response {
		startTime := time.Now()
		response := next(request)
		duration := time.Since(startTime)

		fmt.Printf("Request took %s\n", duration.String())

		return response
	}
}

type IPFilterMiddleware struct {
	AllowedIPs []string
}

func (ipm IPFilterMiddleware) Apply(next AppHandler) AppHandler {
	return func(request *Request) Response {
		clientIP := strings.Split(request.HTTPRequest.RemoteAddr, ":")[0]

		// Check if the client IP is allowed
		if !ipm.isAllowedIP(clientIP) {
			return HtmlResponse{
				Code:    http.StatusForbidden,
				Content: "Access denied",
			}
		}

		return next(request)
	}
}

func (ipm IPFilterMiddleware) isAllowedIP(ip string) bool {
	for _, allowedIP := range ipm.AllowedIPs {
		if allowedIP == ip {
			return true
		}
	}
	return false
}

type RateLimiterMiddleware struct {
	maxRequests   int
	interval      time.Duration
	ipCounters    map[string]int
	ipLastRequest map[string]time.Time
	lock          sync.Mutex
}

// Create in-memory rate limiter
// The limiter is limiting IP to make more then maxRequests in the interval duration.
func NewRateLimiterMiddleware(maxRequests int, interval time.Duration) *RateLimiterMiddleware {
	return &RateLimiterMiddleware{
		maxRequests:   maxRequests,
		interval:      interval,
		ipCounters:    make(map[string]int),
		ipLastRequest: make(map[string]time.Time),
		lock:          sync.Mutex{},
	}
}

func (rlm *RateLimiterMiddleware) Apply(next AppHandler) AppHandler {
	return func(request *Request) Response {
		clientIP := strings.Split(request.HTTPRequest.RemoteAddr, ":")[0]

		resp := func() Response {
			rlm.lock.Lock()
			defer rlm.lock.Unlock()

			// Find last request time
			lastRequestTime, ok := rlm.ipLastRequest[clientIP]
			if !ok {
				lastRequestTime = time.Unix(0, 0)
			}

			// Reset the counter for this client IP if the time interval has elapsed since the last request
			timeSinceLastRequest := time.Since(lastRequestTime)
			if timeSinceLastRequest > rlm.interval {
				rlm.ipCounters[clientIP] = 0
				rlm.ipLastRequest[clientIP] = time.Now()
			}

			// Increase the counter for this client IP
			rlm.ipCounters[clientIP]++

			// If the number of requests from this client IP exceeds the limit, reject the request
			if rlm.ipCounters[clientIP] > rlm.maxRequests {
				return HtmlResponse{
					Code:    http.StatusTooManyRequests,
					Content: "Rate limit exceeded",
				}
			}

			return nil
		}()

		if resp != nil {
			return resp
		}

		// Call the next middleware or handler
		response := next(request)

		return response
	}
}
