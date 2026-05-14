package scraper

import (
	"github.com/FaintLocket424/opengrid-bridge/internal/models"
)

// Scraper defines the capabilities that any timing scraper must have to be valid.
type Scraper interface {
	GetLiveTiming() (*models.RaceResultScrape, error)
	GetPracticeRaceResult(practice, round int) (*models.RaceResultScrape, error)
	GetQualiRaceResult(heat, round int) (*models.RaceResultScrape, error)
	GetFinalRaceResult(final, leg int) (*models.RaceResultScrape, error)
	GetRaceResultsIndex() (*models.RaceResultsIndexScrape, error)
}
