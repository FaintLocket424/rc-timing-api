package service

import (
	"fmt"
	"sync"
	"time"

	"github.com/FaintLocket424/rc-timing-api/internal/models"
	"github.com/FaintLocket424/rc-timing-api/internal/scraper"
	"golang.org/x/sync/singleflight"
)

type EventCache struct {
	Results map[string]*models.RawResult
}

type DataStore struct {
	mu     sync.RWMutex
	events map[string]*EventCache
	group  singleflight.Group
}

func NewStore() *DataStore {
	return &DataStore{events: make(map[string]*EventCache)}
}

func (s *DataStore) getCache(url string) *EventCache {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.events[url] == nil {
		s.events[url] = &EventCache{
			Results: make(map[string]*models.RawResult),
		}
	}
	return s.events[url]
}

func (s *DataStore) GetHeatResult(url string, scrp scraper.Scraper, heat, round int) (models.RawResult, error) {
	cache := s.getCache(url)
	key := fmt.Sprintf("h%dr%d", heat, round)

	s.mu.RLock()
	raw, exists := cache.Results[key]
	s.mu.RUnlock()

	if exists && time.Since(raw.ScrapedAt) < 30*time.Second {
		return *raw, nil
	}

	val, err, _ := s.group.Do(url+key, func() (interface{}, error) {
		fetched, err := scrp.ScrapeResult(url, heat, round)
		if err != nil {
			return nil, err
		}

		s.mu.Lock()
		cache.Results[key] = &fetched
		s.mu.Unlock()

		return fetched, nil
	})

	return val.(models.RawResult), err
}
