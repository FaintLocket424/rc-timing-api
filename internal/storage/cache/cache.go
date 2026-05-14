package cache

import (
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/FaintLocket424/opengrid-bridge/internal/models"
)

var ErrNotFound = errors.New("event data not found")

type EventData struct {
	Live            *models.RaceResultScrape
	PracticeResults map[models.HeatRound]*models.RaceResultScrape
	QualiResults    map[models.HeatRound]*models.RaceResultScrape
	FinalResults    map[models.HeatRound]*models.RaceResultScrape
}

type Cache struct {
	data   map[string]*EventData
	mu     sync.RWMutex
	logger *slog.Logger
}

func NewCache() *Cache {
	c := &Cache{
		data:   make(map[string]*EventData),
		logger: slog.Default().With("component", "cache"),
	}

	go c.startReaper()
	c.logger.Info("Reaper process started")
	return c
}

func (c *Cache) startReaper() {
	for {
		now := time.Now()

		nextMidnight := time.Date(now.Year(), now.Month(), now.Day()+1, 4, 0, 0, 0, now.Location())

		duration := nextMidnight.Sub(now)

		time.Sleep(duration)

		c.mu.Lock()
		c.data = make(map[string]*EventData)
		c.mu.Unlock()

		c.logger.Info("Cache cleared by midnight reaper process")
	}
}

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

func (c *Cache) SaveLiveTiming(url string, model *models.RaceResultScrape) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	ed := c.getOrCreateEvent(url)
	ed.Live = model

	slog.Debug("Saved live timing to cache", "url", url)

	return nil
}

func (c *Cache) GetLiveTiming(url string) (*models.RaceResultScrape, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	ed, ok := c.data[url]
	if !ok || ed.Live == nil {
		return nil, ErrNotFound
	}

	return ed.Live, nil
}

func (c *Cache) SavePracticeRaceResult(url string, model *models.RaceResultScrape) error {
	if model.PracticeNumber == nil || model.Round == nil {
		return fmt.Errorf("cannot store race result, practice=%v; round=%v", model.PracticeNumber, model.Round)
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	ed := c.getOrCreateEvent(url)
	ed.PracticeResults[models.HeatRound{Heat: *model.PracticeNumber, Round: *model.Round}] = model

	slog.Debug("Saved practice race result to cache", "url", url, "heat", *model.PracticeNumber, "round", *model.Round)

	return nil
}

func (c *Cache) GetPracticeRaceResult(url string, heat, round int) (*models.RaceResultScrape, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	key := models.HeatRound{Heat: heat, Round: round}

	ed, ok := c.data[url]
	if !ok || ed.PracticeResults[key] == nil {
		return nil, ErrNotFound
	}

	return ed.PracticeResults[key], nil
}

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

func (c *Cache) GetQualiRaceResult(url string, heat, round int) (*models.RaceResultScrape, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	key := models.HeatRound{Heat: heat, Round: round}

	ed, ok := c.data[url]
	if !ok || ed.QualiResults[key] == nil {
		return nil, ErrNotFound
	}

	return ed.QualiResults[key], nil
}

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

func (c *Cache) GetFinalRaceResult(url string, final, round int) (*models.RaceResultScrape, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	key := models.HeatRound{Heat: final, Round: round}

	ed, ok := c.data[url]
	if !ok || ed.FinalResults[key] == nil {
		return nil, ErrNotFound
	}

	return ed.FinalResults[key], nil
}
