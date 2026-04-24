package manager

import (
	"context"
	"sync"
	"time"

	"github.com/FaintLocket424/rc-timing-api/internal/scraper"
	"github.com/FaintLocket424/rc-timing-api/internal/storage"
)

type WorkerState struct {
	lastAccessed time.Time
	cancel       context.CancelFunc
}

type Manager struct {
	store         storage.Store
	activeWorkers map[string]*WorkerState // Set of URLs which have an active worker
	mu            sync.Mutex
}

func NewManager(store storage.Store) *Manager {
	m := &Manager{
		store:         store,
		activeWorkers: make(map[string]*WorkerState),
	}

	go m.startReaper()

	return m
}

func (m *Manager) EnsureTracking(url string) {
	m.mu.Lock()
	_, ok := m.activeWorkers[url]
	if !ok {
		ctx, cancel := context.WithCancel(context.Background())
		m.activeWorkers[url] = &WorkerState{lastAccessed: time.Now(), cancel: cancel}
		go m.startWorker(ctx, url)
	} else {
		m.activeWorkers[url].lastAccessed = time.Now()
	}
	m.mu.Unlock()
}

func (m *Manager) startWorker(ctx context.Context, url string) {
	s := scraper.NewScraperForURL(url)
	ticker := time.NewTicker(10 * time.Second)

	for {
		select {
		case <-ticker.C:
			model, err := s.GetLiveTiming()
			if err == nil {
				m.store.SaveLiveTiming(url, model)
			}
		case <-ctx.Done():
			ticker.Stop()
			return
		}
	}
}

func (m *Manager) startReaper() {
	ticker := time.NewTicker(1 * time.Minute)

	for {
		m.mu.Lock()

		for url, state := range m.activeWorkers {
			if time.Since(state.lastAccessed) > 30*time.Minute {
				state.cancel()
				delete(m.activeWorkers, url)
			}
		}

		m.mu.Unlock()

		<-ticker.C
	}
}
