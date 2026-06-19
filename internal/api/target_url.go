package api

import (
	"net/http"
	"net/url"
	"strings"

	"codeberg.org/OpenGrid-RC/bridge/internal/api/responses"
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

// ExtractTargetURL is a Gin middleware that enforces the presence of the `target_url`
// query parameter, checks that it is a valid URL, and propagates initialization errors.
func ExtractTargetURL(h *Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		target := c.Query("target_url")

		if !isValidURL(target) {
			responses.RespondError(c, http.StatusBadRequest, "Missing required query parameter: target_url")
			return
		}

		started, err := h.tracker.EnsureTracking(target)
		if err != nil {
			// Check if the error is specifically due to outdated/old timing data.
			if strings.Contains(strings.ToLower(err.Error()), "outdated") ||
				strings.Contains(strings.ToLower(err.Error()), "age of the page") {
				responses.RespondError(c, http.StatusUnprocessableEntity, "The target server has outdated timing data from a previous event")
				return
			}

			// Handle other scraper initialization failures (e.g., connection timed out, DNS lookup failure).
			responses.RespondError(c, http.StatusBadGateway, "Failed to initialize scraper: "+err.Error())
			return
		}

		if started {
			responses.RespondAccepted(c, "Starting tracking event, please poll again in a few seconds")
			return
		}

		c.Set("target_url", target)

		c.Next()
	}
}
