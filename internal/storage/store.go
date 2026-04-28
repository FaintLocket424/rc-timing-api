package storage

import "github.com/FaintLocket424/rc-timing-api/internal/models"

type Store interface {
	SaveLiveTiming(url string, model *models.ResultScrape) error
	GetLiveTiming(url string) (*models.ResultScrape, error)
}
