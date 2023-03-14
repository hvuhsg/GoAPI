package middlewares

import (
	"net/http"
	"time"

	"github.com/hvuhsg/goapi/request"
	"github.com/hvuhsg/goapi/responses"
	"github.com/patrickmn/go-cache"
)

type CacheMiddleware struct {
	cache     *cache.Cache
	keyPrefix string
}

func NewCacheMiddleware(expiration time.Duration, keyPrefix string) *CacheMiddleware {
	return &CacheMiddleware{
		cache:     cache.New(expiration, time.Minute),
		keyPrefix: keyPrefix,
	}
}

func (cm *CacheMiddleware) Apply(next AppHandler) AppHandler {
	return func(request *request.Request) responses.Response {
		// Generate a cache key from the request URL
		cacheKey := cm.keyPrefix + request.HTTPRequest.URL.String()

		// Check if the response is already in the cache
		cachedResponse, found := cm.cache.Get(cacheKey)
		if found {
			response := cachedResponse.(responses.Response)
			// If the response is in the cache, return it directly
			return response
		}

		// Call the next middleware or handler
		response := next(request)

		// If the response status code is 200 OK, cache the response
		if response.StatusCode() == http.StatusOK {
			cm.cache.Set(cacheKey, response, cache.DefaultExpiration)
		}

		return response
	}
}
