// Package scraper handles fetching pages from existing live timing software pages
// and converting it to the internal data representations in the [internal/models]
// package.
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
