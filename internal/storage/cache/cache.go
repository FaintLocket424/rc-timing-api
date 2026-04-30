package cache

import (
	"errors"
	"fmt"
	"sync"

	"github.com/FaintLocket424/rc-timing-api/internal/models"
)

var ErrNotFound = errors.New("event data not found")

type HeatRound struct {
	Heat  int
	Round int
}

type EventData struct {
	Live            *models.RaceResultScrape
	PracticeResults map[HeatRound]*models.RaceResultScrape
	QualiResults    map[HeatRound]*models.RaceResultScrape
	FinalResults    map[HeatRound]*models.RaceResultScrape
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

func (c *Cache) getOrCreateEvent(url string) *EventData {
	if _, ok := c.data[url]; !ok {
		c.data[url] = &EventData{
			PracticeResults: make(map[HeatRound]*models.RaceResultScrape),
			QualiResults:    make(map[HeatRound]*models.RaceResultScrape),
			FinalResults:    make(map[HeatRound]*models.RaceResultScrape),
		}
	}

	return c.data[url]
}

func (c *Cache) SaveLiveTiming(url string, model *models.RaceResultScrape) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	ed := c.getOrCreateEvent(url)
	ed.Live = model

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
	if model.HeatNumber == nil || model.Round == nil {
		return fmt.Errorf("Cannot store race result, heat=%v; round=%v", model.HeatNumber, model.Round)
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	ed := c.getOrCreateEvent(url)
	ed.PracticeResults[HeatRound{*model.HeatNumber, *model.Round}] = model

	return nil
}

func (c *Cache) GetPracticeRaceResult(url string, heat, round int) (*models.RaceResultScrape, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	key := HeatRound{heat, round}

	ed, ok := c.data[url]
	if !ok || ed.PracticeResults[key] == nil {
		return nil, ErrNotFound
	}

	return ed.PracticeResults[key], nil
}

func (c *Cache) SaveQualiRaceResult(url string, model *models.RaceResultScrape) error {
	if model.HeatNumber == nil || model.Round == nil {
		return fmt.Errorf("Cannot store race result, heat=%v; round=%v", model.HeatNumber, model.Round)
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	ed := c.getOrCreateEvent(url)
	ed.QualiResults[HeatRound{*model.HeatNumber, *model.Round}] = model

	return nil
}

func (c *Cache) GetQualiRaceResult(url string, heat, round int) (*models.RaceResultScrape, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	key := HeatRound{heat, round}

	ed, ok := c.data[url]
	if !ok || ed.QualiResults[key] == nil {
		return nil, ErrNotFound
	}

	return ed.QualiResults[key], nil
}

func (c *Cache) SaveFinalRaceResult(url string, model *models.RaceResultScrape) error {
	if model.HeatNumber == nil || model.Round == nil {
		return fmt.Errorf("Cannot store race result, heat=%v; round=%v", model.HeatNumber, model.Round)
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	ed := c.getOrCreateEvent(url)
	ed.FinalResults[HeatRound{*model.HeatNumber, *model.Round}] = model

	return nil
}

func (c *Cache) GetFinalRaceResult(url string, final, round int) (*models.RaceResultScrape, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	key := HeatRound{final, round}

	ed, ok := c.data[url]
	if !ok || ed.FinalResults[key] == nil {
		return nil, ErrNotFound
	}

	return ed.FinalResults[key], nil
}
