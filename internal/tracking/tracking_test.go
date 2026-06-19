package tracking

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/FaintLocket424/opengrid-bridge/internal/models"
	"github.com/FaintLocket424/opengrid-bridge/internal/scraper"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// ptr returns a pointer to the value input.
func ptr[T any](v T) *T {
	return &v
}

type SupervisorTestSuite struct {
	suite.Suite
	originalWorkerLifespan       time.Duration
	originalReaperInterval       time.Duration
	originalLiveTimingInterval   time.Duration
	originalScheduleInterval     time.Duration
	originalResultsIndexInterval time.Duration
	originalResultFetchInterval  time.Duration

	mockStore   *MockStore
	mockFactory *MockScraperFactory
	mockScraper *MockScraper
}

func (suite *SupervisorTestSuite) SetupSuite() {
	// Back up original package configurations
	suite.originalWorkerLifespan = WorkerLifespan
	suite.originalReaperInterval = ReaperInterval
	suite.originalLiveTimingInterval = LiveTimingInterval
	suite.originalScheduleInterval = ScheduleInterval
	suite.originalResultsIndexInterval = ResultsIndexInterval
	suite.originalResultFetchInterval = ResultFetchInterval
}

func (suite *SupervisorTestSuite) TearDownSuite() {
	// Restore original package configurations
	WorkerLifespan = suite.originalWorkerLifespan
	ReaperInterval = suite.originalReaperInterval
	LiveTimingInterval = suite.originalLiveTimingInterval
	ScheduleInterval = suite.originalScheduleInterval
	ResultsIndexInterval = suite.originalResultsIndexInterval
	ResultFetchInterval = suite.originalResultFetchInterval
}

func (suite *SupervisorTestSuite) SetupTest() {
	// Use large/safe values by default to prevent the background reaper
	// daemon from prematurely killing active workers during standard tests.
	WorkerLifespan = 1 * time.Hour
	ReaperInterval = 1 * time.Hour
	LiveTimingInterval = 1 * time.Second
	ScheduleInterval = 1 * time.Second
	ResultsIndexInterval = 1 * time.Second
	ResultFetchInterval = 1 * time.Millisecond

	suite.mockStore = new(MockStore)
	suite.mockFactory = new(MockScraperFactory)
	suite.mockScraper = new(MockScraper)
}

// TestEnsureTracking_Lifecycle verifies starting a tracking worker, checking last accessed
// updates, and shutting it down cleanly.
func (suite *SupervisorTestSuite) TestEnsureTracking_Lifecycle() {
	url := "http://fake.url"

	suite.mockFactory.On("Create", url).Return(suite.mockScraper, nil)
	suite.mockScraper.On("GetLiveTiming").Return(&models.RaceResultScrape{}, nil)
	suite.mockScraper.On("GetRaceSchedule").Return(&models.RaceScheduleScrape{}, nil)
	suite.mockScraper.On("GetRaceResultsIndex").Return(&models.RaceResultsIndexScrape{}, nil)

	suite.mockStore.On("SaveLiveTiming", url, mock.Anything).Return(nil)
	suite.mockStore.On("SaveRaceSchedule", url, mock.Anything).Return(nil)
	suite.mockStore.On("SaveRaceResultsIndex", url, mock.Anything).Return(nil)

	supervisor := NewSupervisor(suite.mockStore, suite.mockFactory)

	// Ensure tracking for the first time starts the worker loop
	started, err := supervisor.EnsureTracking(url)
	suite.Require().NoError(err)
	suite.Require().True(started)

	supervisor.mu.Lock()
	state1, ok := supervisor.activeWorkers[url]
	supervisor.mu.Unlock()
	suite.Require().True(ok)
	suite.Require().NotZero(state1.lastAccessed)

	originalLastAccessed := state1.lastAccessed

	// Second invocation should not spawn a new worker, but rather refresh the existing worker's access time
	time.Sleep(5 * time.Millisecond)
	startedAgain, err := supervisor.EnsureTracking(url)
	suite.Require().NoError(err)
	suite.Require().False(startedAgain)

	supervisor.mu.Lock()
	state2, ok := supervisor.activeWorkers[url]
	supervisor.mu.Unlock()
	suite.Require().True(ok)
	suite.Require().True(state2.lastAccessed.After(originalLastAccessed))

	// Clean up worker
	supervisor.mu.Lock()
	supervisor.reapWorker(url)
	supervisor.mu.Unlock()
}

// TestEnsureTracking_ScraperInitFailure asserts that the supervisor registers the scraper initialization
// failure on the worker state, allowing subsequent EnsureTracking calls to retrieve the cached error.
func (suite *SupervisorTestSuite) TestEnsureTracking_ScraperInitFailure() {
	url := "http://failure.url"

	suite.mockFactory.On("Create", url).Return(nil, errors.New("cannot create scraper"))

	supervisor := NewSupervisor(suite.mockStore, suite.mockFactory)

	started, err := supervisor.EnsureTracking(url)
	suite.Require().NoError(err)
	suite.Require().True(started)

	// Wait asynchronously for the background worker to fail and write the error to the state
	suite.Require().Eventually(func() bool {
		supervisor.mu.Lock()
		defer supervisor.mu.Unlock()
		state, ok := supervisor.activeWorkers[url]
		return ok && state.err != nil
	}, 200*time.Millisecond, 10*time.Millisecond, "Worker state should transition to capture the scraper initialization failure")

	// Ensure that subsequent requests retrieve the cached error
	startedAgain, errAgain := supervisor.EnsureTracking(url)
	suite.Require().False(startedAgain)
	suite.Require().Error(errAgain)
	suite.Require().Contains(errAgain.Error(), "cannot create scraper")
}

// TestReaper_StopsIdleWorkers checks that the automated background reaper shuts down workers
// that have exceeded their lifespans without updates.
func (suite *SupervisorTestSuite) TestReaper_StopsIdleWorkers() {
	url := "http://idle.url"

	// Configure short durations specifically for this test
	WorkerLifespan = 10 * time.Millisecond
	ReaperInterval = 5 * time.Millisecond

	suite.mockFactory.On("Create", url).Return(suite.mockScraper, nil)
	suite.mockScraper.On("GetLiveTiming").Return(&models.RaceResultScrape{}, nil)
	suite.mockScraper.On("GetRaceSchedule").Return(&models.RaceScheduleScrape{}, nil)
	suite.mockScraper.On("GetRaceResultsIndex").Return(&models.RaceResultsIndexScrape{}, nil)

	suite.mockStore.On("SaveLiveTiming", url, mock.Anything).Return(nil)
	suite.mockStore.On("SaveRaceSchedule", url, mock.Anything).Return(nil)
	suite.mockStore.On("SaveRaceResultsIndex", url, mock.Anything).Return(nil)

	supervisor := NewSupervisor(suite.mockStore, suite.mockFactory)

	started, err := supervisor.EnsureTracking(url)
	suite.Require().NoError(err)
	suite.Require().True(started)

	suite.Require().Eventually(func() bool {
		supervisor.mu.Lock()
		defer supervisor.mu.Unlock()
		_, ok := supervisor.activeWorkers[url]
		return !ok
	}, 200*time.Millisecond, 10*time.Millisecond, "Automated reaper did not stop the idle worker")
}

// TestWorker_IncrementalResultScraping tests that cached values are skipped during scraping,
// only querying the network for missing items before saving them.
func (suite *SupervisorTestSuite) TestWorker_IncrementalResultScraping() {
	url := "http://incremental.url"

	indexData := &models.RaceResultsIndexScrape{
		Practice: &models.ResultStatus{
			Results: []models.HeatRound{
				{Heat: 1, Round: 1}, // Cached in Storage
				{Heat: 1, Round: 2}, // Missing in Storage
			},
		},
	}

	result1 := &models.RaceResultScrape{PracticeNumber: ptr(1), Round: ptr(1)}
	result2 := &models.RaceResultScrape{PracticeNumber: ptr(1), Round: ptr(2)}

	suite.mockFactory.On("Create", url).Return(suite.mockScraper, nil)
	suite.mockScraper.On("GetLiveTiming").Return(&models.RaceResultScrape{}, nil)
	suite.mockScraper.On("GetRaceSchedule").Return(&models.RaceScheduleScrape{}, nil)
	suite.mockScraper.On("GetRaceResultsIndex").Return(indexData, nil)

	// Mock Storage cache check: returning successful hit for Round 1, missing for Round 2
	suite.mockStore.On("GetPracticeRaceResult", url, models.HeatRound{Heat: 1, Round: 1}).Return(result1, nil)
	suite.mockStore.On("GetPracticeRaceResult", url, models.HeatRound{Heat: 1, Round: 2}).Return(nil, errors.New("not found"))

	// Expecting only Round 2 to query the scraper and store
	suite.mockScraper.On("GetPracticeRaceResult", models.HeatRound{Heat: 1, Round: 2}).Return(result2, nil)

	suite.mockStore.On("SaveLiveTiming", url, mock.Anything).Return(nil)
	suite.mockStore.On("SaveRaceSchedule", url, mock.Anything).Return(nil)
	suite.mockStore.On("SaveRaceResultsIndex", url, indexData).Return(nil)

	supervisor := NewSupervisor(suite.mockStore, suite.mockFactory)

	// Concurrency synchronization hook to avoid sleep-dependent races
	var wg sync.WaitGroup
	wg.Add(1)

	suite.mockStore.On("SavePracticeRaceResult", url, result2).Run(func(_ mock.Arguments) {
		wg.Done()
	}).Return(nil)

	_, err := supervisor.EnsureTracking(url)
	suite.Require().NoError(err)

	c := make(chan struct{})
	go func() {
		wg.Wait()
		close(c)
	}()

	select {
	case <-c:
		// Succeeded
	case <-time.After(500 * time.Millisecond):
		suite.Fail("Timed out waiting for worker processing loop to execute")
	}

	supervisor.mu.Lock()
	supervisor.reapWorker(url)
	supervisor.mu.Unlock()

	// Ensure Round 1 (cache hit) was never pulled from the scraper
	suite.mockScraper.AssertNotCalled(suite.T(), "GetPracticeRaceResult", models.HeatRound{Heat: 1, Round: 1})
	// Ensure Round 2 (cache miss) was successfully requested
	suite.mockScraper.AssertCalled(suite.T(), "GetPracticeRaceResult", models.HeatRound{Heat: 1, Round: 2})
}

// --- Mocks Definitions ---

type MockStore struct {
	mock.Mock
}

func (m *MockStore) SaveLiveTiming(url string, model *models.RaceResultScrape) error {
	return m.Called(url, model).Error(0)
}

func (m *MockStore) GetLiveTiming(url string) (*models.RaceResultScrape, error) {
	args := m.Called(url)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.RaceResultScrape), args.Error(1)
}

func (m *MockStore) SaveRaceSchedule(url string, model *models.RaceScheduleScrape) error {
	return m.Called(url, model).Error(0)
}

func (m *MockStore) GetRaceSchedule(url string) (*models.RaceScheduleScrape, error) {
	args := m.Called(url)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.RaceScheduleScrape), args.Error(1)
}

func (m *MockStore) SaveRaceResultsIndex(url string, model *models.RaceResultsIndexScrape) error {
	return m.Called(url, model).Error(0)
}

func (m *MockStore) GetRaceResultsIndex(url string) (*models.RaceResultsIndexScrape, error) {
	args := m.Called(url)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.RaceResultsIndexScrape), args.Error(1)
}

func (m *MockStore) GetPracticeRaceResult(url string, hr models.HeatRound) (*models.RaceResultScrape, error) {
	args := m.Called(url, hr)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.RaceResultScrape), args.Error(1)
}

func (m *MockStore) SavePracticeRaceResult(url string, model *models.RaceResultScrape) error {
	return m.Called(url, model).Error(0)
}

func (m *MockStore) GetQualiRaceResult(url string, hr models.HeatRound) (*models.RaceResultScrape, error) {
	args := m.Called(url, hr)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.RaceResultScrape), args.Error(1)
}

func (m *MockStore) SaveQualiRaceResult(url string, model *models.RaceResultScrape) error {
	return m.Called(url, model).Error(0)
}

func (m *MockStore) GetFinalRaceResult(url string, hr models.HeatRound) (*models.RaceResultScrape, error) {
	args := m.Called(url, hr)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.RaceResultScrape), args.Error(1)
}

func (m *MockStore) SaveFinalRaceResult(url string, model *models.RaceResultScrape) error {
	return m.Called(url, model).Error(0)
}

type MockScraperFactory struct {
	mock.Mock
}

func (m *MockScraperFactory) Create(url string) (scraper.Scraper, error) {
	args := m.Called(url)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(scraper.Scraper), args.Error(1)
}

type MockScraper struct {
	mock.Mock
}

func (m *MockScraper) GetLiveTiming() (*models.RaceResultScrape, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.RaceResultScrape), args.Error(1)
}

func (m *MockScraper) GetPracticeRaceResult(hr models.HeatRound) (*models.RaceResultScrape, error) {
	args := m.Called(hr)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.RaceResultScrape), args.Error(1)
}

func (m *MockScraper) GetQualiRaceResult(hr models.HeatRound) (*models.RaceResultScrape, error) {
	args := m.Called(hr)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.RaceResultScrape), args.Error(1)
}

func (m *MockScraper) GetFinalRaceResult(hr models.HeatRound) (*models.RaceResultScrape, error) {
	args := m.Called(hr)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.RaceResultScrape), args.Error(1)
}

func (m *MockScraper) GetRaceResultsIndex() (*models.RaceResultsIndexScrape, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.RaceResultsIndexScrape), args.Error(1)
}

func (m *MockScraper) GetRaceSchedule() (*models.RaceScheduleScrape, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.RaceScheduleScrape), args.Error(1)
}

func TestSupervisorSuite(t *testing.T) {
	suite.Run(t, new(SupervisorTestSuite))
}
