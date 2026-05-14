package manager

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

type Manager struct {
	store         storage.Store
	activeWorkers map[string]*workerState // Set of URLs which have an active worker
	mu            sync.Mutex
	logger        *slog.Logger
}

func NewManager(store storage.Store) *Manager {
	m := &Manager{
		store:         store,
		activeWorkers: make(map[string]*workerState),
		logger:        slog.Default().With("component", "manager"),
	}

	go m.startReaper()
	m.logger.Info("Reaper process started")
	return m
}

func (m *Manager) cleanupWorkerInternal(url string) {
	if state, ok := m.activeWorkers[url]; ok {
		state.cancel()
		delete(m.activeWorkers, url)
		m.logger.Info("worker stopped", "url", url)
	}
}

func (m *Manager) startWorker(ctx context.Context, url string) {
	logger := m.logger.With("url", url)
	s, err := scraper.NewScraperForURL(url)
	if err != nil {
		logger.Error("failed to init scraper", "err", err)
		m.mu.Lock()
		m.cleanupWorkerInternal(url)
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
			process(model.Quali, m.store.GetQualiRaceResult, s.GetQualiRaceResult, m.store.SaveQualiRaceResult, "quali")
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

func (m *Manager) EnsureTracking(url string) (workerStarted bool) {
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

func (m *Manager) startReaper() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		m.mu.Lock()

		for url, state := range m.activeWorkers {
			if time.Since(state.lastAccessed) > workerLifespan {
				m.logger.Info("reaping idle worker", "url", url)
				m.cleanupWorkerInternal(url)
			}
		}

		m.mu.Unlock()
	}
}
