package cache

import (
	"errors"
	"sync"

	"github.com/FaintLocket424/rc-timing-api/internal/models"
)

type EventData struct {
	Live *models.ResultScrape
}

type Cache struct {
	data map[string]*EventData
	mu   sync.RWMutex
}

func NewCache() *Cache {
	return &Cache{
		data: make(map[string]*EventData),
	}
}

func (c *Cache) SaveLiveTiming(url string, model *models.ResultScrape) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.data[url]; !ok {
		c.data[url] = &EventData{}
	}

	c.data[url].Live = model

	return nil
}

func (c *Cache) GetLiveTiming(url string) (*models.ResultScrape, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if _, ok := c.data[url]; !ok {
		return nil, errors.New("Event Data not found in cache for url")
	}

	return c.data[url].Live, nil
}
