// Package tracking handles the lifecycle of the tracking goroutines used to
// scrape each active timing URL.
package tracking

import (
	"context"
	"log/slog"
	"sort"
	"sync"
	"time"

	"github.com/FaintLocket424/opengrid-bridge/internal/models"
	"github.com/FaintLocket424/opengrid-bridge/internal/scraper"
	"github.com/FaintLocket424/opengrid-bridge/internal/storage"
)

var workerLifespan = 30 * time.Minute

type workerState struct {
	lastAccessed time.Time
	cancel       context.CancelFunc
}

// Supervisor holds references to the active store, the scraper factory, a Mutex and a logger.
// It also holds a map of timing URLs to a struct representing the state of
// the tracking goroutine associated with the URL.
// The struct has methods for making sure a URL is tracked, starting a
// new worker goroutine etc.
type Supervisor struct {
	store          storage.Store
	scraperFactory *scraper.Factory
	activeWorkers  map[string]*workerState // Set of URLs which have an active worker
	mu             sync.Mutex
	logger         *slog.Logger
}

// NewSupervisor creates a new supervisor object with a reference to the input store and
// factory scraper. It then starts the reaper goroutine, whose job it is to kill
// inactive worker goroutines.
func NewSupervisor(store storage.Store, scraperFactory *scraper.Factory) *Supervisor {
	m := &Supervisor{
		store:          store,
		scraperFactory: scraperFactory,
		activeWorkers:  make(map[string]*workerState),
		logger:         slog.Default().With("component", "supervisor"),
	}

	go m.startReaper()
	m.logger.Info("Reaper process started")
	return m
}

// reapWorker kills the goroutine associated with the url parameter, giving it
// the cancel signal and deleting it from the active workers struct.
func (m *Supervisor) reapWorker(url string) {
	if state, ok := m.activeWorkers[url]; ok {
		state.cancel()
		delete(m.activeWorkers, url)
		m.logger.Info("worker stopped", "url", url)
	}
}

// startWorker is the main function of the tracking goroutines.
// It handles creating the scraper for the URL and using it to
// scrape the url on a fixed interval.
func (m *Supervisor) startWorker(ctx context.Context, url string) {
	logger := m.logger.With("url", url)
	s, err := m.scraperFactory.Create(url)
	if err != nil {
		logger.Error("failed to init scraper", "err", err)
		m.mu.Lock()
		m.reapWorker(url)
		m.mu.Unlock()
		return
	}

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		reqCount := 0
		maxReqs := 5

		canRequest := func() bool {
			if reqCount >= maxReqs {
				return false
			}
			reqCount++
			return true
		}

		if canRequest() {
			if model, err := s.GetLiveTiming(); err != nil {
				logger.Warn("live timing scrape failed", "err", err)
			} else if err := m.store.SaveLiveTiming(url, model); err != nil {
				logger.Error("live timing storage failed", "err", err)
			}
		}

		if model, err := s.GetRaceResultsIndex(); err == nil {
			logger.Debug("Scraped race results index successfully")
			process := func(
				results map[models.HeatRound]struct{},
				checker func(string, int, int) (*models.RaceResultScrape, error),
				fetcher func(int, int) (*models.RaceResultScrape, error),
				save func(string, *models.RaceResultScrape) error,
				kind string,
			) {
				type item struct{ heat, round int }
				var list []item
				for key := range results {
					list = append(list, item{heat: key.Heat, round: key.Round})
				}

				// Sort by Round then Heat
				sort.Slice(list, func(i, j int) bool {
					if list[i].round != list[j].round {
						return list[i].round < list[j].round
					}
					return list[i].heat < list[j].heat
				})

				for _, entry := range list {
					// Check cache first
					if _, err := checker(url, entry.heat, entry.round); err != nil {
						// Check rate limit before network call
						if !canRequest() {
							return
						}

						time.Sleep(500 * time.Millisecond)

						res, err := fetcher(entry.heat, entry.round)
						if err != nil {
							logger.Warn(kind+" scrape failed", "heat", entry.heat, "round", entry.round, "err", err)
							continue
						}
						if err := save(url, res); err != nil {
							logger.Error(kind+" storage failed", "err", err)
						}
					}
				}
			}

			// Execute processing with limit
			process(model.Practice, m.store.GetPracticeRaceResult, s.GetPracticeRaceResult, m.store.SavePracticeRaceResult, "practice")
			process(model.Qualifying, m.store.GetQualiRaceResult, s.GetQualiRaceResult, m.store.SaveQualiRaceResult, "quali")
			process(model.Finals, m.store.GetFinalRaceResult, s.GetFinalRaceResult, m.store.SaveFinalRaceResult, "final")
		} else {
			logger.Warn("results index scrape failed", "err", err)
		}

		select {
		case <-ticker.C:
			logger.Debug("worker waking up")
			continue
		case <-ctx.Done():
			logger.Info("worker shutting down")
			return
		}
	}
}

// EnsureTracking creates a tracking goroutine for the input URL, if it doesn't
// exist already. It also handles updating the last accessed time for the URL
// so the supervisor knows to keep it alive.
func (m *Supervisor) EnsureTracking(url string) (workerStarted bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if state, ok := m.activeWorkers[url]; ok {
		state.lastAccessed = time.Now()
		return false
	}

	ctx, cancel := context.WithCancel(context.Background())
	m.activeWorkers[url] = &workerState{lastAccessed: time.Now(), cancel: cancel}
	go m.startWorker(ctx, url)
	m.logger.Info("Worker started", "url", url)
	return true
}

// startReaper is the main function of the reaper goroutine, which runs on a
// fixed interval and reaps any worker tracking goroutines that have not been
// accessed for longer than the workerLifespan.
func (m *Supervisor) startReaper() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		m.mu.Lock()

		for url, state := range m.activeWorkers {
			if time.Since(state.lastAccessed) > workerLifespan {
				m.logger.Info("reaping idle worker", "url", url)
				m.reapWorker(url)
			}
		}

		m.mu.Unlock()
	}
}
