// Package storage holds the Store interface, representing a data structure for storing previously
// fetched data for quick retrieval later.
package storage

import "codeberg.org/OpenGrid-RC/bridge/internal/models"

// Store represents a data cache with methods to save and retrieve timing data.
type Store interface {
	SaveLiveTiming(url string, model *models.RaceResultScrape) error
	GetLiveTiming(url string) (*models.RaceResultScrape, error)

	SavePracticeRaceResult(url string, model *models.RaceResultScrape) error
	GetPracticeRaceResult(url string, hr models.HeatRound) (*models.RaceResultScrape, error)

	SaveQualiRaceResult(url string, model *models.RaceResultScrape) error
	GetQualiRaceResult(url string, hr models.HeatRound) (*models.RaceResultScrape, error)

	SaveFinalRaceResult(url string, model *models.RaceResultScrape) error
	GetFinalRaceResult(url string, hr models.HeatRound) (*models.RaceResultScrape, error)

	SaveRaceResultsIndex(url string, model *models.RaceResultsIndexScrape) error
	GetRaceResultsIndex(url string) (*models.RaceResultsIndexScrape, error)

	SaveRaceSchedule(url string, model *models.RaceScheduleScrape) error
	GetRaceSchedule(url string) (*models.RaceScheduleScrape, error)
}
