package service

import (
	"strings"
	"sync"
	"time"

	"github.com/FaintLocket424/rc-timing-api/internal/models"
	"golang.org/x/sync/singleflight"
)

type DataStore struct {
	mu    sync.RWMutex
	cache map[string]models.EventReport
	group singleflight.Group // Prevents duplicate scrapers
}

func NewStore() *DataStore {
	return &DataStore{
		cache: make(map[string]models.EventReport),
	}
}

func (s *DataStore) GetEvent(targetURL string, scraper func(string) (models.EventReport, error)) (models.EventReport, error) {
	s.mu.RLock()
	data, exists := s.cache[targetURL]
	s.mu.RUnlock()

	if exists && !data.IsExpired(1*time.Minute) {
		return data, nil
	}

	val, err, _ := s.group.Do(targetURL, func() (interface{}, error) {
		newData, err := scraper(targetURL)
		if err != nil {
			return nil, err
		}

		s.mu.Lock()
		s.cache[targetURL] = newData
		s.mu.Unlock()

		return newData, nil
	})

	if err != nil {
		return models.EventReport{}, err
	}

	return val.(models.EventReport), nil
}

func (s *DataStore) GetSchedule(targetURL string, scraper func(string) (models.EventReport, error)) ([]models.HeatSummary, error) {
	event, err := s.GetEvent(targetURL, scraper)
	if err != nil {
		return nil, err
	}

	var summary []models.HeatSummary
	for _, h := range event.Heats {
		summary = append(summary, h.HeatSummary)
	}
	return summary, nil
}

func (s *DataStore) GetHeatsByType(url, sessionType string, scraper func(string) (models.EventReport, error)) ([]models.Heat, error) {
	event, err := s.GetEvent(url, scraper)
	if err != nil {
		return nil, err
	}

	var filtered []models.Heat
	for _, h := range event.Heats {
		if strings.EqualFold(h.Type, sessionType) {
			filtered = append(filtered, h)
		}
	}
	return filtered, nil
}
