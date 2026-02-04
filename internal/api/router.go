package api

import (
	"fmt"
	"net/http"

	"github.com/FaintLocket424/rc-timing-api/internal/api/handlers"
	"github.com/FaintLocket424/rc-timing-api/internal/service"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(store *service.DataStore) *gin.Engine {
	r := gin.Default()

	r.Use(gzip.Gzip(gzip.DefaultCompression))

	r.Use(func(c *gin.Context) {
		c.Next()

		if updated, exists := c.Get("last_updated"); exists {
			etag := fmt.Sprintf("%v", updated)
			if c.GetHeader("If-None-Match") == etag {
				c.AbortWithStatus(http.StatusNotModified)
			}
			c.Header("ETag", etag)
		}
	})

	eventHandler := handlers.NewEventHandler(store)

	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("/api/v1")
	{
		v1.GET("/liv", eventHandler.GetLiveRace)
	}

	return r
}
