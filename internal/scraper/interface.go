package scraper

import "github.com/FaintLocket424/rc-timing-api/internal/models"

// Scraper defines the capabilities that any timing scraper must have to be valid.
type Scraper interface {
	GetLiveTiming() (*models.ResultScrape, error)
}
