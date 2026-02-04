package scraper

import (
	"github.com/FaintLocket424/rc-timing-api/internal/models"
)

type BBKScraper struct{}

func (s *BBKScraper) ScrapeQualifyingResult(baseURL string, heat, round int) (models.CachedHeatResult, error) {
	return models.CachedHeatResult{}, nil
}
