package cache

import (
	"errors"

	"github.com/FaintLocket424/rc-timing-api/internal/models"
)

type EventData struct {
}

type Cache struct {
	data map[string]EventData
}

func NewCache() *Cache {
	return &Cache{
		make(map[string]EventData),
	}
}

func (c Cache) SaveLiveTiming(url string, model *models.LiveTimingScrape) error {
	return errors.New("Not implemented")
}

func (c Cache) GetLiveTiming(url string) (*models.LiveTimingScrape, error) {
	return nil, errors.New("Not implemented")
}
