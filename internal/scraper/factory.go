package scraper

import (
	"fmt"
	"net/http"

	"github.com/FaintLocket424/rc-timing-api/internal/scraper/bbk"
	"github.com/PuerkitoBio/goquery"
)

// NewScraperForURL takes in a target url and detects which Scraper is suitable.
func NewScraperForURL(url string) (Scraper, error) {
	client := http.DefaultClient

	res, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code error: %d, %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)

	author, exists := doc.Find("meta[name='author']").Attr("content")

	if exists && author == "bbkRClive" {
		return &bbk.BBKScraper{
			Target: url,
			Client: client,
		}, nil
	}

	return nil, fmt.Errorf("Unable to determine scraper for target: \"%s\"", url)
}
