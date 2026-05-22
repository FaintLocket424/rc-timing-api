package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRateLimiter(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// Limit 1 request per second, burst size of 2
	r.Use(RateLimiter(1, 2))
	r.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	performRequest := func(ip string) int {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("X-Forwarded-For", ip)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		return w.Code
	}

	// Client 1: First 2 requests (burst) should pass
	assert.Equal(t, http.StatusOK, performRequest("192.168.1.1"))
	assert.Equal(t, http.StatusOK, performRequest("192.168.1.1"))

	// Client 1: 3rd request should hit the limit (429)
	assert.Equal(t, http.StatusTooManyRequests, performRequest("192.168.1.1"))

	// Client 2: Should have its own bucket, so it should pass
	assert.Equal(t, http.StatusOK, performRequest("192.168.1.2"))
}
