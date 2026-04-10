package scraper

import "github.com/FaintLocket424/rc-timing-api/internal/scraper/bbk"

// NewScraperForURL takes in a target url and detects which Scraper is suitable.
func NewScraperForURL(url string) Scraper {
	// Placeholder that always returns bbk
	return &bbk.BBKScraper{}
}
