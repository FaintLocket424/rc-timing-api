package manager

import (
	"context"
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
}

func NewManager(store storage.Store) *Manager {
	m := &Manager{
		store:         store,
		activeWorkers: make(map[string]*workerState),
	}

	go m.startReaper()

	return m
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

func (m *Manager) startWorker(ctx context.Context, url string) {
	s := scraper.NewScraperForURL(url)
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	if model, err := s.GetLiveTiming(); err == nil {
		m.store.SaveLiveTiming(url, model)
	}

	for {
		select {
		case <-ticker.C:
			if model, err := s.GetLiveTiming(); err == nil {
				m.store.SaveLiveTiming(url, model)
			}
		case <-ctx.Done():
			return
		}
	}
}

func (m *Manager) startReaper() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		m.mu.Lock()

		for url, state := range m.activeWorkers {
			if time.Since(state.lastAccessed) > 30*time.Minute {
				state.cancel()
				delete(m.activeWorkers, url)
			}
		}

		m.mu.Unlock()
	}
}
