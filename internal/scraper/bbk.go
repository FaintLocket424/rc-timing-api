package scraper

import (
	"github.com/FaintLocket424/rc-timing-api/internal/models"
)

type BBKScraper struct{}

func (s *BBKScraper) ScrapeResult(baseURL string, heat, round int) (models.RawResult, error) {
	return models.RawResult{}, nil
}
