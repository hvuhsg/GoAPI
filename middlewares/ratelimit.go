package middlewares

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/hvuhsg/goapi/request"
	"github.com/hvuhsg/goapi/responses"
)

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
	return func(request *request.Request) responses.Response {
		clientIP := strings.Split(request.HTTPRequest.RemoteAddr, ":")[0]

		resp := func() responses.Response {
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
				return responses.NewHTMLResponse("Rate limit exceeded", http.StatusTooManyRequests)
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
