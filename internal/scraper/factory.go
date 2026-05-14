package scraper

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/FaintLocket424/opengrid-bridge/internal/scraper/bbk"
	"github.com/PuerkitoBio/goquery"
)

type RealHTTPClient struct {
	client *http.Client
}

func (c *RealHTTPClient) Get(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36")
	return c.client.Do(req)
}

// NewScraperForURL takes in a target url and detects which Scraper is suitable.
func NewScraperForURL(url string) (Scraper, error) {
	client := &RealHTTPClient{http.DefaultClient}

	res, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to reach index, status code error: %d, %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)

	author, exists := doc.Find("meta[name='author']").Attr("content")

	if exists && author == "bbkRClive" {
		return &bbk.BBKScraper{
			Target: url,
			Client: client,
		}, nil
	}

	return nil, errors.New("unable to determine scraper")
}
