package storage

import "github.com/FaintLocket424/rc-timing-api/internal/models"

type Store interface {
	SaveLiveTiming(url string, model *models.LiveTimingScrape) error
	GetLiveTiming(url string) (*models.LiveTimingScrape, error)
}
