package main

import (
	"fmt"
	"net/http"

	"github.com/FaintLocket424/rc-timing-api/internal/scraper"
	"github.com/FaintLocket424/rc-timing-api/internal/service"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func main() {
	store := service.NewStore()
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

	api := r.Group("/api/v1")
	{
		api.GET("/schedule", func(c *gin.Context) {
			data, err := store.GetSchedule(c.Query("url"), scraper.Scrape)

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.Set("last_updated", data[0].ScheduledAt.Unix())
			c.JSON(http.StatusOK, data)
		})

		api.GET("/heats/:type", func(c *gin.Context) {
			data, _ := store.GetHeatsByType(c.Query("url"), c.Param("type"), scraper.Scrape)
			c.JSON(http.StatusOK, data)
		})

		api.GET("/live", func(c *gin.Context) {
			event, _ := store.GetEvent(c.Query("url"), scraper.Scrape)

			for _, h := range event.Heats {
				if h.Status == "Live" {
					c.JSON(http.StatusOK, h)
					return
				}
			}

			c.JSON(http.StatusNotFound, gin.H{"message": "No live race"})
		})
	}

	err := r.Run(":8080")
	if err != nil {
		println(err.Error())
		return
	}
}
