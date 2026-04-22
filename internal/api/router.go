package api

import (
	"net/http"

	"github.com/FaintLocket424/rc-timing-api/internal/storage/cache"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.SetTrustedProxies(nil)
	handler := NewHandler(cache.NewCache())

	v1 := r.Group("/api/v1")
	{
		v1.GET("/ping", func(c *gin.Context) {
			c.String(http.StatusOK, "pong")
		})

		v1.GET("/live", handler.GetLiveTiming)
	}

	return r
}
