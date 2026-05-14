package scraper

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/FaintLocket424/opengrid-bridge/internal/scraper/bbk"
	"github.com/PuerkitoBio/goquery"
)

// NewScraperForURL takes in a target url and detects which Scraper is suitable.
func NewScraperForURL(url string) (scraper Scraper, err error) {
	client := NewClient()

	res, fetchErr := client.Get(url)
	if fetchErr != nil {
		return nil, fetchErr
	}

	defer func() {
		if closeErr := res.Body.Close(); closeErr != nil {
			err = errors.Join(err, closeErr)
		}
	}()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to reach index, status code error: %d, %s", res.StatusCode, res.Status)
	}

	doc, docErr := goquery.NewDocumentFromReader(res.Body)
	if docErr != nil {
		return nil, docErr
	}

	author, exists := doc.Find("meta[name='author']").Attr("content")

	if exists && author == "bbkRClive" {
		return &bbk.Scraper{
			Target: url,
			Client: client,
		}, nil
	}

	return nil, errors.New("unable to determine scraper")
}
