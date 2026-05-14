package scraper

import (
	"net/http"
	"time"
)

const userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36"

// const openGridUserAgent = "OpenGrid Timing Bridge/1.0 (+https://github.com/FaintLocket424/opengrid-bridge)"

type userAgentTransport struct {
	transport http.RoundTripper
}

// RoundTrip modifies http requests to add a user agent header to each one.
func (t *userAgentTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	reqClone := req.Clone(req.Context())

	reqClone.Header.Set("User-Agent", userAgent)

	return t.transport.RoundTrip(reqClone)
}

// NewClient creates a new HTTP Client that uses our custom userAgentTransport struct
// to customise the requests before being sent. Plus, it adds a 10 second timeout
// on requests.
func NewClient() *http.Client {
	return &http.Client{
		Transport: &userAgentTransport{
			transport: http.DefaultTransport,
		},
		Timeout: 10 * time.Second,
	}
}
