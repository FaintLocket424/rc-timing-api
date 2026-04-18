package bbk

import (
	"net/http"
)

type HTTPClient interface {
	Get(url string) (*http.Response, error)
}

type BBKScraper struct {
	Target string
	Client HTTPClient
}
