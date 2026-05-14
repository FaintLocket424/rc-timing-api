package scraper

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/FaintLocket424/opengrid-bridge/internal/scraper/bbk"
	"github.com/PuerkitoBio/goquery"
)

// Factory is a scraper factory with a Create method that makes new scrapers.
// It holds a reference to the programVersion to be used in the scraper's
// User-Agent http header.
type Factory struct {
	client *http.Client
}

// NewFactory creates a new Scraper Factory with the program version injected.
func NewFactory(programVersion string) *Factory {
	return &Factory{
		client: NewClient(programVersion),
	}
}

// Create makes a new scraper using the factory.
func (f *Factory) Create(url string) (scraper Scraper, err error) {
	client := f.client

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
		return bbk.NewScraper(url, client), nil
	}

	return nil, errors.New("unable to determine scraper")
}
