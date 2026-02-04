package scraper

import "github.com/FaintLocket424/rc-timing-api/internal/models"

type Scraper interface {
	ScrapePracticeResult(baseURL string, heat, round int) (models.CachedHeatResult, error)
	ScrapeQualifyingResult(baseURL string, heat, round int) (models.CachedHeatResult, error)
}

func GetScraper(software string) Scraper {
	switch software {
	//case "bbk":
	//	return &BBKScraper{}
	case "fake":
		return &FakeScraper{}
	default:
		return nil
	}
}
