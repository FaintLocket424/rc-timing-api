package scraper

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

type userAgentTransport struct {
	transport http.RoundTripper
	userAgent string
}

// RoundTrip modifies http requests to add a user agent header to each one.
// As well as adding accepted data formats and language to better circumvent
// simple anti-bot firewalls.
func (t *userAgentTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	reqClone := req.Clone(req.Context())

	reqClone.Header.Set("User-Agent", t.userAgent)
	reqClone.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	reqClone.Header.Set("Accept-Language", "en-US,en;q=0.5")

	return t.transport.RoundTrip(reqClone)
}

// NewClient creates a new HTTP Client that uses our custom userAgentTransport struct
// to customise the requests before being sent. Plus, it adds a 10 second timeout
// on requests.
func NewClient(programVersion, programCommit string) *http.Client {
	fullVersion := programVersion
	if programCommit != "" && programCommit != "unknown" && programCommit != "dev" {
		fullVersion = fmt.Sprintf("%s-%s", programVersion, programCommit)
	}

	userAgent := strings.Join([]string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64)",
		"AppleWebKit/537.36 (KHTML, like Gecko)",
		"Chrome/148.0.7778.98",
		"Safari/537.36",
		fmt.Sprintf("OpenGridBridge/%s (+https://github.com/FaintLocket424/opengrid-bridge)", fullVersion),
	}, " ")

	return &http.Client{
		Transport: &userAgentTransport{
			transport: http.DefaultTransport,
			userAgent: userAgent,
		},
		Timeout: 10 * time.Second,
	}
}
