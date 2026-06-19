// Package cache is an in-memory cache that follows the [internal/storage] Storage interface.
package cache

import (
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"codeberg.org/OpenGrid-RC/bridge/internal/models"
)

var errNotFound = errors.New("event data not found")

// EventData holds all the data which has been scraped from an event so far.
type EventData struct {
	Live            *models.RaceResultScrape
	ResultsIndex    *models.RaceResultsIndexScrape
	RaceSchedule    *models.RaceScheduleScrape
	PracticeResults map[models.HeatRound]*models.RaceResultScrape
	QualiResults    map[models.HeatRound]*models.RaceResultScrape
	FinalResults    map[models.HeatRound]*models.RaceResultScrape
}

// Cache is an in-memory cache that implements the [internal/storage] Storage interface.
// It maps timing URL links to an EventData struct, as well as holds a mutex for
// concurrent access.
type Cache struct {
	data   map[string]*EventData
	mu     sync.RWMutex
	logger *slog.Logger
}

// NewCache creates a new cache object and starts the cache reaper goroutine.
func NewCache() *Cache {
	c := &Cache{
		data:   make(map[string]*EventData),
		logger: slog.Default().With("component", "cache"),
	}

	go c.startReaper()
	c.logger.Info("Reaper process started")
	return c
}

// startReaper is the main function of the cache reaper goroutine.
// It sleeps until 4am and then cleans out the cache since events rarely
// go for multiple days and the cache can be refilled later.
func (c *Cache) startReaper() {
	for {
		now := time.Now()

		nextFourAM := time.Date(now.Year(), now.Month(), now.Day()+1, 4, 0, 0, 0, now.Location())

		duration := nextFourAM.Sub(now)

		time.Sleep(duration)

		c.mu.Lock()
		c.data = make(map[string]*EventData)
		c.mu.Unlock()

		c.logger.Info("Cache cleared by 4am reaper process")
	}
}

// getOrCreateEvent retrieves the EventData for a URL from the cache, or creates a
// new empty object.
func (c *Cache) getOrCreateEvent(url string) *EventData {
	if _, ok := c.data[url]; !ok {
		c.data[url] = &EventData{
			PracticeResults: make(map[models.HeatRound]*models.RaceResultScrape),
			QualiResults:    make(map[models.HeatRound]*models.RaceResultScrape),
			FinalResults:    make(map[models.HeatRound]*models.RaceResultScrape),
		}
	}

	return c.data[url]
}

// SaveLiveTiming saves a live race result scrape for a URL into the cache.
func (c *Cache) SaveLiveTiming(url string, model *models.RaceResultScrape) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	ed := c.getOrCreateEvent(url)
	ed.Live = model

	slog.Debug("Saved live timing to cache", "url", url)

	return nil
}

// GetLiveTiming retrieves the currently stored live race scrape for a URL stored in the cache.
func (c *Cache) GetLiveTiming(url string) (*models.RaceResultScrape, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	ed, ok := c.data[url]
	if !ok || ed.Live == nil {
		return nil, errNotFound
	}

	return ed.Live, nil
}

// SavePracticeRaceResult saves a scraped practice result for a URL into the cache.
// It uses the data in the scrape to work out the practice heat and round it represents.
func (c *Cache) SavePracticeRaceResult(url string, model *models.RaceResultScrape) error {
	if model.PracticeNumber == nil || model.Round == nil {
		return fmt.Errorf("cannot store race result, practice=%v; round=%v",
			model.PracticeNumber, model.Round)
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	ed := c.getOrCreateEvent(url)
	ed.PracticeResults[models.HeatRound{Heat: *model.PracticeNumber, Round: *model.Round}] = model

	slog.Debug("Saved practice race result to cache", "url", url,
		"heat", *model.PracticeNumber, "round", *model.Round)

	return nil
}

// GetPracticeRaceResult retrieves the currently stored practice race result from a given heat/round
// for a URL in the cache.
func (c *Cache) GetPracticeRaceResult(url string, hr models.HeatRound) (*models.RaceResultScrape, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	ed, ok := c.data[url]
	if !ok || ed.PracticeResults[hr] == nil {
		return nil, errNotFound
	}

	return ed.PracticeResults[hr], nil
}

// SaveQualiRaceResult saves a scraped qualifying result for a URL into the cache.
// It uses the data in the scrape to work out the qualifying heat and round it represents.
func (c *Cache) SaveQualiRaceResult(url string, model *models.RaceResultScrape) error {
	if model.HeatNumber == nil || model.Round == nil {
		return fmt.Errorf("cannot store race result, heat=%v; round=%v", model.HeatNumber, model.Round)
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	ed := c.getOrCreateEvent(url)
	ed.QualiResults[models.HeatRound{Heat: *model.HeatNumber, Round: *model.Round}] = model

	slog.Debug("Saved quali race result to cache", "url", url, "heat", *model.HeatNumber, "round", *model.Round)

	return nil
}

// GetQualiRaceResult retrieves the currently stored qualifying race result from a given heat/round
// for a URL in the cache.
func (c *Cache) GetQualiRaceResult(url string, hr models.HeatRound) (*models.RaceResultScrape, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	ed, ok := c.data[url]
	if !ok || ed.QualiResults[hr] == nil {
		return nil, errNotFound
	}

	return ed.QualiResults[hr], nil
}

// SaveFinalRaceResult saves a scraped final result for a URL into the cache.
// It uses the data in the scrape to work out the final number and round it represents.
func (c *Cache) SaveFinalRaceResult(url string, model *models.RaceResultScrape) error {
	if model.FinalNumber == nil || model.Round == nil {
		return fmt.Errorf("cannot store race result, final=%v; round=%v", model.FinalNumber, model.Round)
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	ed := c.getOrCreateEvent(url)
	ed.FinalResults[models.HeatRound{Heat: *model.FinalNumber, Round: *model.Round}] = model

	slog.Debug("Saved final race result to cache", "url", url, "heat", *model.FinalNumber, "round", *model.Round)

	return nil
}

// GetFinalRaceResult retrieves the currently stored final race result from a given final/round
// for a URL in the cache.
func (c *Cache) GetFinalRaceResult(url string, hr models.HeatRound) (*models.RaceResultScrape, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	ed, ok := c.data[url]
	if !ok || ed.FinalResults[hr] == nil {
		return nil, errNotFound
	}

	return ed.FinalResults[hr], nil
}

// SaveRaceResultsIndex saves a scraped results index for a URL into the cache.
func (c *Cache) SaveRaceResultsIndex(url string, model *models.RaceResultsIndexScrape) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	ed := c.getOrCreateEvent(url)
	ed.ResultsIndex = model

	slog.Debug("Saved results index to cache", "url", url)

	return nil
}

// GetRaceResultsIndex retrieves the currently stored results index for a URL in the cache.
func (c *Cache) GetRaceResultsIndex(url string) (*models.RaceResultsIndexScrape, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	ed, ok := c.data[url]
	if !ok || ed.ResultsIndex == nil {
		return nil, errNotFound
	}

	return ed.ResultsIndex, nil
}

// SaveRaceSchedule saves a scraped race schedule for a URL into the cache.
func (c *Cache) SaveRaceSchedule(url string, model *models.RaceScheduleScrape) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	ed := c.getOrCreateEvent(url)
	ed.RaceSchedule = model

	slog.Debug("Saved race schedule to cache", "url", url)

	return nil
}

// GetRaceSchedule retrieves the currently stored race schedule for a URL in the cache.
func (c *Cache) GetRaceSchedule(url string) (*models.RaceScheduleScrape, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	ed, ok := c.data[url]
	if !ok || ed.RaceSchedule == nil {
		return nil, errNotFound
	}

	return ed.RaceSchedule, nil
}
