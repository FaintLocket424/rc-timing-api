package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"codeberg.org/OpenGrid-RC/bridge/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

const targetURL = "http://example.com"

// MockStore acts as a fake cache that the test can "store" data in.
type MockStore struct {
	mock.Mock
}

// SaveLiveTiming is a mock method.
func (m *MockStore) SaveLiveTiming(url string, model *models.RaceResultScrape) error {
	args := m.Called(url, model)
	return args.Error(0)
}

// GetLiveTiming is a mock method.
func (m *MockStore) GetLiveTiming(url string) (*models.RaceResultScrape, error) {
	args := m.Called(url)
	var res *models.RaceResultScrape
	if args.Get(0) != nil {
		res = args.Get(0).(*models.RaceResultScrape)
	}
	return res, args.Error(1)
}

// SaveRaceSchedule is a mock method.
func (m *MockStore) SaveRaceSchedule(url string, model *models.RaceScheduleScrape) error {
	args := m.Called(url, model)
	return args.Error(0)
}

// GetRaceSchedule is a mock method.
func (m *MockStore) GetRaceSchedule(url string) (*models.RaceScheduleScrape, error) {
	args := m.Called(url)
	var res *models.RaceScheduleScrape
	if args.Get(0) != nil {
		res = args.Get(0).(*models.RaceScheduleScrape)
	}
	return res, args.Error(1)
}

// SavePracticeRaceResult is a mock method.
func (m *MockStore) SavePracticeRaceResult(url string, model *models.RaceResultScrape) error {
	args := m.Called(url, model)
	return args.Error(0)
}

// GetPracticeRaceResult is a mock method.
func (m *MockStore) GetPracticeRaceResult(url string, hr models.HeatRound) (*models.RaceResultScrape, error) {
	args := m.Called(url, hr)

	var res *models.RaceResultScrape
	if args.Get(0) != nil {
		res = args.Get(0).(*models.RaceResultScrape)
	}
	return res, args.Error(1)
}

// SaveQualiRaceResult is a mock method.
func (m *MockStore) SaveQualiRaceResult(url string, model *models.RaceResultScrape) error {
	args := m.Called(url, model)
	return args.Error(0)
}

// GetQualiRaceResult is a mock method.
func (m *MockStore) GetQualiRaceResult(url string, hr models.HeatRound) (*models.RaceResultScrape, error) {
	args := m.Called(url, hr)
	var res *models.RaceResultScrape
	if args.Get(0) != nil {
		res = args.Get(0).(*models.RaceResultScrape)
	}
	return res, args.Error(1)
}

// SaveFinalRaceResult is a mock method.
func (m *MockStore) SaveFinalRaceResult(url string, model *models.RaceResultScrape) error {
	args := m.Called(url, model)
	return args.Error(0)
}

// GetFinalRaceResult is a mock method.
func (m *MockStore) GetFinalRaceResult(url string, hr models.HeatRound) (*models.RaceResultScrape, error) {
	args := m.Called(url, hr)
	var res *models.RaceResultScrape
	if args.Get(0) != nil {
		res = args.Get(0).(*models.RaceResultScrape)
	}
	return res, args.Error(1)
}

// SaveRaceResultsIndex is a mock method.
func (m *MockStore) SaveRaceResultsIndex(url string, model *models.RaceResultsIndexScrape) error {
	args := m.Called(url, model)
	return args.Error(0)
}

// GetRaceResultsIndex is a mock method.
func (m *MockStore) GetRaceResultsIndex(url string) (*models.RaceResultsIndexScrape, error) {
	args := m.Called(url)
	var res *models.RaceResultsIndexScrape
	if args.Get(0) != nil {
		res = args.Get(0).(*models.RaceResultsIndexScrape)
	}
	return res, args.Error(1)
}

// MockTracker represents a fake tracker.
type MockTracker struct {
	mock.Mock
}

// EnsureTracking is a mock method returning two parameters to conform with the updated design.
func (m *MockTracker) EnsureTracking(url string) (bool, error) {
	args := m.Called(url)
	return args.Bool(0), args.Error(1)
}

// HandlerTestSuite is a testify test suite for testing the Handler struct.
type HandlerTestSuite struct {
	suite.Suite
	mockStore   *MockStore
	mockTracker *MockTracker
	router      *gin.Engine
}

// SetupTest initialises the test suite, creating a mock store and tracker.
func (suite *HandlerTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)

	suite.mockStore = new(MockStore)
	suite.mockTracker = new(MockTracker)

	suite.router = SetupRouter(suite.mockStore, suite.mockTracker, "test", "commit")
}

// TearDownTest asserts that the mock store and tracker returned all
// their values.
func (suite *HandlerTestSuite) TearDownTest() {
	suite.mockStore.AssertExpectations(suite.T())
	suite.mockTracker.AssertExpectations(suite.T())
}

// =========================
// Live Timing Tests
// =========================

// TestLive_MissingURL tests that the API returns a StatusBadRequest 400
// if no target_url query param is passed in.
func (suite *HandlerTestSuite) TestLive_MissingURL() {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/live", nil)
	suite.router.ServeHTTP(w, req)

	suite.Require().Equal(http.StatusBadRequest, w.Code)
	suite.JSONEq(`{
		"status": "error",
		"message": "Missing required query parameter: target_url"
	}`, w.Body.String())
}

// TestLive_TrackerInitializing tests that the api returns the message
// "Starting tracking event, please poll again in a few seconds" when EnsureTracking
// returns true and the cache returns an error.
func (suite *HandlerTestSuite) TestLive_TrackerInitializing() {
	suite.mockTracker.On("EnsureTracking", targetURL).Return(true, nil).Once()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/live?target_url="+targetURL, nil)
	suite.router.ServeHTTP(w, req)

	suite.Require().Equal(http.StatusAccepted, w.Code)
	suite.JSONEq(`{
		"status": "success",
		"message": "Starting tracking event, please poll again in a few seconds"
	}`, w.Body.String())
}

// TestLive_TrackerOutdatedData tests that the API returns StatusUnprocessableEntity 422
// with a descriptive error message when outdated timing files are detected.
func (suite *HandlerTestSuite) TestLive_TrackerOutdatedData() {
	suite.mockTracker.On("EnsureTracking", targetURL).Return(false, errors.New("timing file is outdated")).Once()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/live?target_url="+targetURL, nil)
	suite.router.ServeHTTP(w, req)

	suite.Require().Equal(http.StatusUnprocessableEntity, w.Code)
	suite.JSONEq(`{
		"status": "error",
		"message": "The target server has outdated timing data from a previous event"
	}`, w.Body.String())
}

// TestLive_TrackerGenericInitError tests that the API returns StatusBadGateway 502
// when a non-outdated initialization error occurs.
func (suite *HandlerTestSuite) TestLive_TrackerGenericInitError() {
	suite.mockTracker.On("EnsureTracking", targetURL).Return(false, errors.New("dns lookup failed")).Once()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/live?target_url="+targetURL, nil)
	suite.router.ServeHTTP(w, req)

	suite.Require().Equal(http.StatusBadGateway, w.Code)
	suite.JSONEq(`{
		"status": "error",
		"message": "Failed to initialize scraper: dns lookup failed"
	}`, w.Body.String())
}

func (suite *HandlerTestSuite) TestLive_Success() {
	className := "2 Wheel Drive"
	dummyData := &models.RaceResultScrape{ClassName: &className}

	suite.mockTracker.On("EnsureTracking", targetURL).Return(false, nil).Once()
	suite.mockStore.On("GetLiveTiming", targetURL).Return(dummyData, nil).Once()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/live?target_url="+targetURL, nil)
	suite.router.ServeHTTP(w, req)

	suite.Require().Equal(http.StatusOK, w.Code)
	suite.JSONEq(`{
		"status": "success",
		"data": {
			"class_name": "2 Wheel Drive"
		}
	}`, w.Body.String())
}

// =========================
// Race Schedule Tests
// =========================

func (suite *HandlerTestSuite) TestSchedule_MissingURL() {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/schedule", nil)
	suite.router.ServeHTTP(w, req)

	suite.Require().Equal(http.StatusBadRequest, w.Code)
	suite.JSONEq(`{
		"status": "error",
		"message": "Missing required query parameter: target_url"
	}`, w.Body.String())
}

func (suite *HandlerTestSuite) TestSchedule_TrackerInitializing() {
	suite.mockTracker.On("EnsureTracking", targetURL).Return(true, nil).Once()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/schedule?target_url="+targetURL, nil)
	suite.router.ServeHTTP(w, req)

	suite.Require().Equal(http.StatusAccepted, w.Code)
	suite.JSONEq(`{
		"status": "success",
		"message": "Starting tracking event, please poll again in a few seconds"
	}`, w.Body.String())
}

func (suite *HandlerTestSuite) TestSchedule_StoreError() {
	suite.mockTracker.On("EnsureTracking", targetURL).Return(false, nil).Once()
	suite.mockStore.On("GetRaceSchedule", targetURL).Return(nil, errors.New("database failure")).Once()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/schedule?target_url="+targetURL, nil)
	suite.router.ServeHTTP(w, req)

	suite.Require().Equal(http.StatusInternalServerError, w.Code)
	suite.JSONEq(`{
		"status": "error",
		"message": "database failure"
	}`, w.Body.String())
}

func (suite *HandlerTestSuite) TestSchedule_Success() {
	title := "Weekend Schedule"
	dummyData := &models.RaceScheduleScrape{Title: &title}

	suite.mockTracker.On("EnsureTracking", targetURL).Return(false, nil).Once()
	suite.mockStore.On("GetRaceSchedule", targetURL).Return(dummyData, nil).Once()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/schedule?target_url="+targetURL, nil)
	suite.router.ServeHTTP(w, req)

	suite.Require().Equal(http.StatusOK, w.Code)
	suite.Contains(w.Body.String(), "Weekend Schedule")
}

// =========================
// Practice Heats Tests
// =========================

// TestPractice_MissingURL tests that the API returns a StatusBadRequest 400
// if no target_url query param is passed in.
func (suite *HandlerTestSuite) TestPractice_MissingURL() {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/results/practice/round/1/heat/2", nil)
	suite.router.ServeHTTP(w, req)

	suite.Require().Equal(http.StatusBadRequest, w.Code)
	suite.JSONEq(`{
		"status": "error",
		"message": "Missing required query parameter: target_url"
	}`, w.Body.String())
}

// TestPractice_TrackerInitializing tests that the api returns the message
// "Starting tracking event, please poll again in a few seconds" when EnsureTracking
// returns true and the cache returns an error.
func (suite *HandlerTestSuite) TestPractice_TrackerInitializing() {
	suite.mockTracker.On("EnsureTracking", targetURL).Return(true, nil).Once()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/results/practice/round/1/heat/2?target_url="+targetURL, nil)
	suite.router.ServeHTTP(w, req)

	suite.Require().Equal(http.StatusAccepted, w.Code)
	suite.JSONEq(`{
		"status": "success",
		"message": "Starting tracking event, please poll again in a few seconds"
	}`, w.Body.String())
}

// TestPractice_StoreError tests that the server returns a code 500
// StatusInternalServerError if the store fails to retrieve the data.
func (suite *HandlerTestSuite) TestPractice_StoreError() {
	suite.mockTracker.On("EnsureTracking", targetURL).Return(false, nil).Once()
	suite.mockStore.On("GetPracticeRaceResult", targetURL, models.HeatRound{Heat: 2, Round: 1}).Return(nil, errors.New("db error")).Once()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/results/practice/round/1/heat/2?target_url="+targetURL, nil)
	suite.router.ServeHTTP(w, req)

	suite.Require().Equal(http.StatusInternalServerError, w.Code)
	suite.JSONEq(`{
		"status": "error",
		"message": "db error"
	}`, w.Body.String())
}

// TestPractice_Success tests if the handler correctly returns the data found
// in the cache.
func (suite *HandlerTestSuite) TestPractice_Success() {
	className := "LMP2"
	suite.mockTracker.On("EnsureTracking", targetURL).Return(false, nil).Once()
	suite.mockStore.On("GetPracticeRaceResult", targetURL, models.HeatRound{Heat: 2, Round: 1}).Return(&models.RaceResultScrape{ClassName: &className}, nil).Once()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/results/practice/round/1/heat/2?target_url="+targetURL, nil)
	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)
	suite.Contains(w.Body.String(), "LMP2")
}

// =========================
// Qualifying Heats Tests
// =========================

func (suite *HandlerTestSuite) TestQuali_MissingURL() {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/results/qualifying/round/3/heat/4", nil)
	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusBadRequest, w.Code)
	suite.JSONEq(`{
		"status": "error",
		"message": "Missing required query parameter: target_url"
	}`, w.Body.String())
}

func (suite *HandlerTestSuite) TestQuali_TrackerInitializing() {
	suite.mockTracker.On("EnsureTracking", targetURL).Return(true, nil).Once()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/results/qualifying/round/3/heat/4?target_url="+targetURL, nil)
	suite.router.ServeHTTP(w, req)

	suite.Require().Equal(http.StatusAccepted, w.Code)
	suite.JSONEq(`{
		"status": "success",
		"message": "Starting tracking event, please poll again in a few seconds"
	}`, w.Body.String())
}

func (suite *HandlerTestSuite) TestQuali_StoreError() {
	suite.mockTracker.On("EnsureTracking", targetURL).Return(false, nil).Once()
	suite.mockStore.On("GetQualiRaceResult", targetURL, models.HeatRound{Heat: 4, Round: 3}).Return(nil, errors.New("timeout")).Once()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/results/qualifying/round/3/heat/4?target_url="+targetURL, nil)
	suite.router.ServeHTTP(w, req)

	suite.Require().Equal(http.StatusInternalServerError, w.Code)
	suite.JSONEq(`{
		"status": "error",
		"message": "timeout"
	}`, w.Body.String())
}

func (suite *HandlerTestSuite) TestQuali_Success() {
	className := "GT4"
	suite.mockTracker.On("EnsureTracking", targetURL).Return(false, nil).Once()
	suite.mockStore.On("GetQualiRaceResult", targetURL, models.HeatRound{Heat: 4, Round: 3}).Return(&models.RaceResultScrape{ClassName: &className}, nil).Once()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/results/qualifying/round/3/heat/4?target_url="+targetURL, nil)
	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)
	suite.Contains(w.Body.String(), "GT4")
}

// =========================
// Final Races Tests
// =========================

func (suite *HandlerTestSuite) TestFinal_MissingURL() {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/results/finals/round/5/final/1", nil)
	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusBadRequest, w.Code)
	suite.JSONEq(`{
		"status": "error",
		"message": "Missing required query parameter: target_url"
	}`, w.Body.String())
}

func (suite *HandlerTestSuite) TestFinal_TrackerInitializing() {
	suite.mockTracker.On("EnsureTracking", targetURL).Return(true, nil).Once()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/results/finals/round/5/final/1?target_url="+targetURL, nil)
	suite.router.ServeHTTP(w, req)

	suite.Require().Equal(http.StatusAccepted, w.Code)
	suite.JSONEq(`{
		"status": "success",
		"message": "Starting tracking event, please poll again in a few seconds"
	}`, w.Body.String())
}

func (suite *HandlerTestSuite) TestFinal_StoreError() {
	suite.mockTracker.On("EnsureTracking", targetURL).Return(false, nil).Once()
	suite.mockStore.On("GetFinalRaceResult", targetURL, models.HeatRound{Heat: 1, Round: 5}).Return(nil, errors.New("db disconnect")).Once()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/results/finals/round/5/final/1?target_url="+targetURL, nil)
	suite.router.ServeHTTP(w, req)

	suite.Require().Equal(http.StatusInternalServerError, w.Code)
	suite.JSONEq(`{
		"status": "error",
		"message": "db disconnect"
	}`, w.Body.String())
}

func (suite *HandlerTestSuite) TestFinal_Success() {
	className := "Hypercar"
	suite.mockTracker.On("EnsureTracking", targetURL).Return(false, nil).Once()
	suite.mockStore.On("GetFinalRaceResult", targetURL, models.HeatRound{Heat: 1, Round: 5}).Return(&models.RaceResultScrape{ClassName: &className}, nil).Once()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/results/finals/round/5/final/1?target_url="+targetURL, nil)
	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)
	suite.Contains(w.Body.String(), "Hypercar")
}

// =========================
// Run the suite
// =========================

// TestHandlerTestSuite runs the Handler test suite.
func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}
