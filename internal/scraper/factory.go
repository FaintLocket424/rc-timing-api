package scraper

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/FaintLocket424/opengrid-bridge/internal/scraper/bbk"
	"github.com/PuerkitoBio/goquery"
)

// Factory creates the correct Scraper for the input URL.
// It holds a reference to the programVersion to be used in the scraper's
// User-Agent http header.
type Factory struct {
	client *http.Client
}

// NewFactory creates a Scraper Factory with the program version injected
// into the singleton HTTP Client used for all networking requests.
func NewFactory(programVersion, programCommit string) *Factory {
	return &Factory{
		client: NewClient(programVersion, programCommit),
	}
}

// Create fetches the index page from the target URL and determines the
// correct scraper for the live timing software being used.
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

	if author, exists := doc.Find("meta[name='author']").Attr("content"); exists && author == "bbkRClive" {
		lastModifiedStr := res.Header.Get("Last-Modified")
		if lastModifiedStr == "" {
			return nil, errors.New("unable to determine age of the page: Last-Modified header is missing")
		}

		lastModified, parseErr := http.ParseTime(lastModifiedStr)
		if parseErr != nil {
			return nil, fmt.Errorf("failed to parse Last-Modified header: %w", parseErr)
		}

		nowUTC := time.Now().UTC()
		lastModifiedUTC := lastModified.UTC()

		y1, m1, d1 := lastModifiedUTC.Date()
		y2, m2, d2 := nowUTC.Date()

		if y1 != y2 || m1 != m2 || d1 != d2 {
			return nil, fmt.Errorf("timing file is outdated: last modified %s, expected current day %s (UTC)",
				lastModifiedUTC.Format("2006-01-02"), nowUTC.Format("2006-01-02"))
		}

		return bbk.NewScraper(url, client), nil
	}

	return nil, errors.New("unable to determine scraper")
}
