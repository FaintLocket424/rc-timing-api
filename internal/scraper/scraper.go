package scraper

import "github.com/FaintLocket424/rc-timing-api/internal/models"

type Scraper interface {
	//ScrapeQualifyingMeta(baseURL string) (models.QualifyingMetaDTO, error)
	//ScrapeQualifyingOveralls(baseURL string, roundID int) ([]models.QualifyingOverallEntryDTO, error)
	ScrapeResult(baseURL string, heat, round int) (models.RawResult, error)
	//ScrapeQualifyingList(baseURL string) ([]models.Heat, error)
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
