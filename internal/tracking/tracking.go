// Package tracking handles the lifecycle of the tracking goroutines used to
// scrape each active timing URL.
package tracking

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/FaintLocket424/opengrid-bridge/internal/models"
	"github.com/FaintLocket424/opengrid-bridge/internal/scraper"
	"github.com/FaintLocket424/opengrid-bridge/internal/storage"
)

// ScraperFactory represents a struct that can create the correct scraper
// for a given URL.
type ScraperFactory interface {
	Create(url string) (scraper.Scraper, error)
}

// Configurable timing intervals.
var (
	WorkerLifespan       = 30 * time.Minute
	ReaperInterval       = 1 * time.Minute
	LiveTimingInterval   = 10 * time.Second
	ScheduleInterval     = 1 * time.Minute
	ResultsIndexInterval = 30 * time.Second
	ResultFetchInterval  = 500 * time.Millisecond
)

type workerState struct {
	lastAccessed time.Time
	cancel       context.CancelFunc
	err          error     // Captures initialization failures
	failedAt     time.Time // Timestamp of the initialization failure
}

// Supervisor holds references to the active store, the scraper factory, a Mutex and a logger.
// It also holds a map of timing URLs to a struct representing the state of
// the tracking goroutine associated with the URL.
type Supervisor struct {
	store          storage.Store
	scraperFactory ScraperFactory
	activeWorkers  map[string]*workerState // Set of URLs which have an active worker
	mu             sync.Mutex
	logger         *slog.Logger
}

// NewSupervisor creates a new supervisor object with a reference to the input store and
// factory scraper. It then starts the reaper goroutine, whose job it is to kill
// inactive worker goroutines.
func NewSupervisor(store storage.Store, scraperFactory ScraperFactory) *Supervisor {
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
// It spawns separate, concurrent loops for each scraper task.
func (m *Supervisor) startWorker(ctx context.Context, url string) {
	logger := m.logger.With("url", url)
	s, err := m.scraperFactory.Create(url)
	if err != nil {
		logger.Error("failed to init scraper", "err", err)
		m.mu.Lock()
		if state, ok := m.activeWorkers[url]; ok {
			state.err = err
			state.failedAt = time.Now()
		}
		m.mu.Unlock()
		return
	}

	var wg sync.WaitGroup

	// 1. Live Timing Scrape Loop
	wg.Add(1)
	go func() {
		defer wg.Done()
		m.scrapeLiveTiming(ctx, s, url, logger)

		ticker := time.NewTicker(LiveTimingInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				m.scrapeLiveTiming(ctx, s, url, logger)
			}
		}
	}()

	// 2. Race Schedule Scrape Loop
	wg.Add(1)
	go func() {
		defer wg.Done()
		m.scrapeRaceSchedule(ctx, s, url, logger)

		ticker := time.NewTicker(ScheduleInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				m.scrapeRaceSchedule(ctx, s, url, logger)
			}
		}
	}()

	// 3. Results Index & Detailed Results Loop
	wg.Add(1)
	go func() {
		defer wg.Done()
		m.scrapeResultsIndex(ctx, s, url, logger)

		ticker := time.NewTicker(ResultsIndexInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				m.scrapeResultsIndex(ctx, s, url, logger)
			}
		}
	}()

	// Wait for all sub-goroutines to finish before exiting
	wg.Wait()
	logger.Info("worker shutting down")
}

func (m *Supervisor) scrapeLiveTiming(ctx context.Context, s scraper.Scraper, url string, logger *slog.Logger) {
	if ctx.Err() != nil {
		return
	}
	model, err := s.GetLiveTiming()
	if err != nil {
		logger.Error("live timing scrape failed", "err", err)
		return
	}
	if err := m.store.SaveLiveTiming(url, model); err != nil {
		logger.Error("live timing storage failed", "err", err)
		return
	}
	logger.Debug("Scraped live timing successfully")
}

func (m *Supervisor) scrapeRaceSchedule(ctx context.Context, s scraper.Scraper, url string, logger *slog.Logger) {
	if ctx.Err() != nil {
		return
	}
	model, err := s.GetRaceSchedule()
	if err != nil {
		logger.Error("race schedule scrape failed", "err", err)
		return
	}
	if err := m.store.SaveRaceSchedule(url, model); err != nil {
		logger.Error("race schedule storage failed", "err", err)
		return
	}
	logger.Debug("Scraped race schedule successfully")
}

func (m *Supervisor) scrapeResultsIndex(ctx context.Context, s scraper.Scraper, url string, logger *slog.Logger) {
	if ctx.Err() != nil {
		return
	}
	model, err := s.GetRaceResultsIndex()
	if err != nil {
		logger.Warn("results index scrape failed", "err", err)
		return
	}
	logger.Debug("Scraped race results index successfully")

	if err := m.store.SaveRaceResultsIndex(url, model); err != nil {
		logger.Error("race result index storage failed", "err", err)
	}

	reqCount := 0
	maxReqs := 5

	canRequest := func() bool {
		if reqCount >= maxReqs {
			return false
		}
		reqCount++
		return true
	}

	process := func(
		results []models.HeatRound,
		checker func(url string, hr models.HeatRound) (*models.RaceResultScrape, error),
		fetcher func(hr models.HeatRound) (*models.RaceResultScrape, error),
		save func(string, *models.RaceResultScrape) error,
		kind string,
	) {
		if len(results) == 0 {
			return
		}

		for _, entry := range results {
			if ctx.Err() != nil {
				return
			}

			// Check cache first
			if _, err := checker(url, entry); err != nil {
				if !canRequest() {
					return
				}

				// Context-aware sleep to avoid blocking worker shutdown during long scrapes
				timer := time.NewTimer(ResultFetchInterval)
				select {
				case <-ctx.Done():
					timer.Stop()
					return
				case <-timer.C:
				}

				res, err := fetcher(entry)
				if err != nil {
					logger.Warn(kind+" scrape failed", "heat", entry.Heat, "round", entry.Round, "err", err)
					continue
				}
				if err := save(url, res); err != nil {
					logger.Error(kind+" storage failed", "err", err)
				}
			}
		}
	}

	// Safe pointer extraction to prevent nil pointer panics
	var practiceResults []models.HeatRound
	if model.Practice != nil {
		practiceResults = model.Practice.Results
	}

	var qualifyingResults []models.HeatRound
	if model.Qualifying != nil {
		qualifyingResults = model.Qualifying.Results
	}

	var finalsResults []models.HeatRound
	if model.Finals != nil {
		finalsResults = model.Finals.Results
	}

	process(practiceResults, m.store.GetPracticeRaceResult, s.GetPracticeRaceResult, m.store.SavePracticeRaceResult, "practice")
	process(qualifyingResults, m.store.GetQualiRaceResult, s.GetQualiRaceResult, m.store.SaveQualiRaceResult, "quali")
	process(finalsResults, m.store.GetFinalRaceResult, s.GetFinalRaceResult, m.store.SaveFinalRaceResult, "final")
}

// EnsureTracking creates a tracking goroutine for the input URL, if it doesn't
// exist already. It handles updating the last accessed time for the URL so the
// supervisor knows to keep it alive. If the background scraper initialization failed,
// it returns the initialization error.
func (m *Supervisor) EnsureTracking(url string) (workerStarted bool, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if state, ok := m.activeWorkers[url]; ok {
		if state.err != nil {
			// If the error occurred more than 1 minute ago, clear the state and allow
			// a retry. Otherwise, continue returning the cached error to prevent hammering.
			if time.Since(state.failedAt) > 1*time.Minute {
				delete(m.activeWorkers, url)
			} else {
				state.lastAccessed = time.Now()
				return false, state.err
			}
		} else {
			state.lastAccessed = time.Now()
			return false, nil
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	m.activeWorkers[url] = &workerState{
		lastAccessed: time.Now(),
		cancel:       cancel,
	}
	go m.startWorker(ctx, url)
	m.logger.Info("Worker started", "url", url)
	return true, nil
}

// startReaper is the main function of the reaper goroutine, which runs on a
// fixed interval and reaps any worker tracking goroutines that have not been
// accessed for longer than the WorkerLifespan.
func (m *Supervisor) startReaper() {
	ticker := time.NewTicker(ReaperInterval)
	defer ticker.Stop()

	for range ticker.C {
		m.mu.Lock()

		for url, state := range m.activeWorkers {
			if time.Since(state.lastAccessed) > WorkerLifespan {
				m.logger.Info("reaping idle worker", "url", url)
				m.reapWorker(url)
			}
		}

		m.mu.Unlock()
	}
}
