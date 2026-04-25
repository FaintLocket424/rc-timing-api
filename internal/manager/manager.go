package manager

import (
	"context"
	"log"
	"log/slog"
	"sync"
	"time"

	"github.com/FaintLocket424/rc-timing-api/internal/scraper"
	"github.com/FaintLocket424/rc-timing-api/internal/storage"
)

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
	s, err := scraper.NewScraperForURL(url)
	if err != nil {
		m.logger.Error("failed to init scraper", "url", url, "err", err)
		m.mu.Lock()
		m.cleanupWorkerInternal(url)
		m.mu.Unlock()
		return
	}

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		model, err := s.GetLiveTiming()
		if err != nil {
			m.logger.Warn("scrape failed", "url", url, "err", err)
		} else if err := m.store.SaveLiveTiming(url, model); err != nil {
			m.logger.Error("storage failed", "url", url, "err", err)
		}

		select {
		case <-ticker.C:
			continue
		case <-ctx.Done():
			log.Printf("worker for %s shutting down", url)
			return
		}
	}
}

func (m *Manager) EnsureTracking(url string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if state, ok := m.activeWorkers[url]; ok {
		state.lastAccessed = time.Now()
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	m.activeWorkers[url] = &workerState{lastAccessed: time.Now(), cancel: cancel}
	go m.startWorker(ctx, url)
}

func (m *Manager) startReaper() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		m.mu.Lock()

		for url, state := range m.activeWorkers {
			if time.Since(state.lastAccessed) > 30*time.Minute {
				m.logger.Info("reaping idle worker", "url", url)
				m.cleanupWorkerInternal(url)
			}
		}

		m.mu.Unlock()
	}
}
