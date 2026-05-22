package middleware

import (
	"net/http"
	"net/url"

	"github.com/FaintLocket424/opengrid-bridge/internal/api/responses"
	"github.com/gin-gonic/gin"
)

// isValidURL checks if a string is a well-formed absolute HTTP/HTTPS URL.
func isValidURL(toTest string) bool {
	u, err := url.ParseRequestURI(toTest)
	if err != nil {
		return false
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return false
	}
	if u.Host == "" {
		return false
	}
	return true
}

// ExtractTargetURL is a Gin middleware that enfoces the presence of the `target_url`
// query parameter, and checks it's a valid URL and returns an error if it's not.
func ExtractTargetURL(c *gin.Context) {
	target := c.Query("target_url")

	if !isValidURL(target) {
		responses.RespondError(c, http.StatusBadRequest, "Missing required query parameter: target_url")
		return
	}

	c.Set("target_url", target)

	c.Next()
}
