// Package scraper handles fetching pages from existing live timing software pages
// and converting it to the internal data representations in the [internal/models]
// package.
package scraper

import (
	"codeberg.org/OpenGrid-RC/bridge/internal/models"
)

// Scraper defines the capabilities that any timing scraper must have to be valid.
type Scraper interface {
	GetLiveTiming() (*models.RaceResultScrape, error)
	GetPracticeRaceResult(models.HeatRound) (*models.RaceResultScrape, error)
	GetQualiRaceResult(models.HeatRound) (*models.RaceResultScrape, error)
	GetFinalRaceResult(models.HeatRound) (*models.RaceResultScrape, error)
	GetRaceResultsIndex() (*models.RaceResultsIndexScrape, error)
	GetRaceSchedule() (*models.RaceScheduleScrape, error)
}
