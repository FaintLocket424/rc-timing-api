package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NotModifiedMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if updated, exists := c.Get("last_updated"); exists {
			etag := fmt.Sprintf("%v", updated)
			if c.GetHeader("If-None-Match") == etag {
				c.AbortWithStatus(http.StatusNotModified)
			}
			c.Header("ETag", etag)
		}
	}
}
