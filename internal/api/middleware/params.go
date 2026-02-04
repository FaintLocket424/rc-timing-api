package middleware

import (
	"net/http"

	"github.com/FaintLocket424/rc-timing-api/internal/scraper"
	"github.com/gin-gonic/gin"
)

const (
	ContextURLKey      = "scrape_url"
	ContextSoftwareKey = "scrape_software"

	QueryParamURL      = "url"
	QueryParamSoftware = "software"
)

func ScrapeParamsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		targetURL := c.Query(QueryParamURL)
		software := c.Query(QueryParamSoftware)

		if targetURL == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required 'url' parameter"})
			c.Abort()
			return
		}

		if software == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required 'software' parameter"})
			c.Abort()
			return
		}

		c.Set(ContextURLKey, targetURL)
		c.Set(ContextSoftwareKey, software)

		c.Next()
	}
}

func GetURL(c *gin.Context) string {
	return c.MustGet(ContextURLKey).(string)
}

func GetSoftware(c *gin.Context) string {
	return c.MustGet(ContextSoftwareKey).(string)
}

func GetScraper(c *gin.Context) scraper.Scraper {
	return scraper.GetScraper(GetSoftware(c))
}
