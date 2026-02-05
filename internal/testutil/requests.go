package testutil

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

func init() {
	// Set GIN to test mode to stop logging
	gin.SetMode(gin.TestMode)
}

func GenerateValidRequest(path, targetURL, software string) *http.Request {
	u, _ := url.Parse(path)
	q := u.Query()
	if targetURL != "" {
		q.Set("url", targetURL)
	}

	if software != "" {
		q.Set("software", software)
	}
	u.RawQuery = q.Encode()

	req, _ := http.NewRequest("GET", u.String(), nil)

	return req
}
