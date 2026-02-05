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
	Meta    *models.CachedMeta
	Results map[string]*models.CachedHeatResult
}

type DataStore struct {
	mu     sync.RWMutex
	events map[string]*EventCache
	group  singleflight.Group
}

func NewStore() *DataStore {
	return &DataStore{events: make(map[string]*EventCache)}
}

// getCache retrieves the current cached data for a given URL.
// A new cache object is created if missing.
func (s *DataStore) getCache(url string) *EventCache {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.events[url] == nil {
		s.events[url] = &EventCache{
			Results: make(map[string]*models.CachedHeatResult),
		}
	}

	return s.events[url]
}

// GetQualiHeatResult probes the data store to see if a valid cache exists.
// If not, it asks the scraper to go fetch it.
func (s *DataStore) GetQualiHeatResult(url string, scraper scraper.Scraper, heat, round int) (models.CachedHeatResult, error) {
	cache := s.getCache(url)
	key := fmt.Sprintf("h%dr%d", heat, round)

	// Apply read lock to the cache, and fetch the cached data.
	s.mu.RLock()
	raw, exists := cache.Results[key]
	s.mu.RUnlock()

	// If the cache is valid, return it.
	if exists && time.Since(raw.ScrapedAt) < 16*time.Second {
		return *raw, nil
	}

	// The cache is invalid, so it needs regenerating.
	// This creates a singleflight block for url+key
	// If multiple requests ask for the same data, only one scrape will occur,
	// with all waiting requests receiving the same data.
	val, err, _ := s.group.Do(url+key, func() (interface{}, error) {
		fetched, err := scraper.ScrapeQualifyingResult(url, heat, round)
		if err != nil {
			return nil, err
		}

		// Write the fetched data to the cache, with a write lock.
		s.mu.Lock()
		cache.Results[key] = &fetched
		s.mu.Unlock()

		return fetched, nil
	})

	return val.(models.CachedHeatResult), err
}

// GetPracticeHeatResult probes the data store to see if a valid cache exists.
// If not, it asks the scraper to go fetch it.
func (s *DataStore) GetPracticeHeatResult(url string, scraper scraper.Scraper, heat, round int) (models.CachedHeatResult, error) {
	cache := s.getCache(url)
	key := fmt.Sprintf("p%dr%d", heat, round)

	// Apply read lock to the cache, and fetch the cached data.
	s.mu.RLock()
	raw, exists := cache.Results[key]
	s.mu.RUnlock()

	// If the cache is valid, return it.
	if exists && time.Since(raw.ScrapedAt) < 16*time.Second {
		return *raw, nil
	}

	// The cache is invalid, so it needs regenerating.
	// This creates a singleflight block for url+key
	// If multiple requests ask for the same data, only one scrape will occur,
	// with all waiting requests receiving the same data.
	val, err, _ := s.group.Do(url+key, func() (interface{}, error) {
		fetched, err := scraper.ScrapePracticeResult(url, heat, round)
		if err != nil {
			return nil, err
		}

		// Write the fetched data to the cache, with a write lock.
		s.mu.Lock()
		cache.Results[key] = &fetched
		s.mu.Unlock()

		return fetched, nil
	})

	return val.(models.CachedHeatResult), err
}

func (s *DataStore) GetEventMeta(url string, scraper scraper.Scraper) (models.CachedMeta, error) {
	cache := s.getCache(url)

	s.mu.RLock()
	raw := cache.Meta
	s.mu.RUnlock()

	if raw != nil && time.Since(cache.Meta.ScrapedAt) < 10*time.Minute {
		return *raw, nil
	}

	val, err, _ := s.group.Do(url+"meta", func() (interface{}, error) {
		fetched, err := scraper.ScrapeEventMeta(url)
		if err != nil {
			return nil, err
		}

		s.mu.Lock()
		cache.Meta = &fetched
		s.mu.Unlock()

		return fetched, nil
	})

	return val.(models.CachedMeta), err
}
