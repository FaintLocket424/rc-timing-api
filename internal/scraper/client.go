package scraper

import (
	"net/http"
	"time"
)

const userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36"

type userAgentTransport struct {
	transport http.RoundTripper
}

func (t *userAgentTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	reqClone := req.Clone(req.Context())

	reqClone.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36")

	return t.transport.RoundTrip(reqClone)
}

func NewClient() *http.Client {
	return &http.Client{
		Transport: &userAgentTransport{
			transport: http.DefaultTransport,
		},
		Timeout: 10 * time.Second,
	}
}
