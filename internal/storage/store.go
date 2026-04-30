package storage

import "github.com/FaintLocket424/rc-timing-api/internal/models"

type Store interface {
	SaveLiveTiming(url string, model *models.RaceResultScrape) error
	GetLiveTiming(url string) (*models.RaceResultScrape, error)

	SavePracticeRaceResult(url string, model *models.RaceResultScrape) error
	GetPracticeRaceResult(url string, heat, round int) (*models.RaceResultScrape, error)

	SaveQualiRaceResult(url string, model *models.RaceResultScrape) error
	GetQualiRaceResult(url string, heat, round int) (*models.RaceResultScrape, error)

	SaveFinalRaceResult(url string, model *models.RaceResultScrape) error
	GetFinalRaceResult(url string, final, round int) (*models.RaceResultScrape, error)
}
