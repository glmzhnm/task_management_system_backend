package middlewares

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type rateLimiter struct {
	mu       sync.Mutex
	requests map[string][]time.Time
	limit    int
	window   time.Duration
}

var limiter = &rateLimiter{
	requests: make(map[string][]time.Time),
	limit:    100,
	window:   time.Minute,
}

func RateLimiterMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		limiter.mu.Lock()
		defer limiter.mu.Unlock()

		now := time.Now()
		windowStart := now.Add(-limiter.window)

		var valid []time.Time
		for _, t := range limiter.requests[ip] {
			if t.After(windowStart) {
				valid = append(valid, t)
			}
		}
		limiter.requests[ip] = valid

		if len(limiter.requests[ip]) >= limiter.limit {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":       "Rate limit exceeded",
				"retry_after": "60s",
			})
			c.Abort()
			return
		}

		limiter.requests[ip] = append(limiter.requests[ip], now)
		c.Next()
	}
}
