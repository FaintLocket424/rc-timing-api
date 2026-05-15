// Package middleware contains Gin middleware functions to be used
// by [internal/api] when creating the router.
package middleware

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type limiterStore struct {
	mu      sync.Mutex
	clients map[string]*rate.Limiter
	r       rate.Limit
	b       int
}

func (s *limiterStore) getLimiter(ip string) *rate.Limiter {
	s.mu.Lock()
	defer s.mu.Unlock()
	lim, exists := s.clients[ip]
	if !exists {
		lim = rate.NewLimiter(s.r, s.b)
		s.clients[ip] = lim
	}
	return lim
}

// RateLimiter handles denying requests from IP addresses that have made too
// many requests. It retrives the client IP address of the request and checks
// if the request is allowed to continue, or it returns a TooManyRequests 429 status.
func RateLimiter(r rate.Limit, b int) gin.HandlerFunc {
	store := limiterStore{
		clients: make(map[string]*rate.Limiter),
		r:       r,
		b:       b,
	}

	return func(c *gin.Context) {
		ip := c.ClientIP()
		lim := store.getLimiter(ip)
		if !lim.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "rate limit exceeded",
			})
			return
		}
		c.Next()
	}
}
