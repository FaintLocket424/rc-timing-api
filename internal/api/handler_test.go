package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/FaintLocket424/opengrid-bridge/internal/models"
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

// SavePracticeRaceResult is a mock method.
func (m *MockStore) SavePracticeRaceResult(url string, model *models.RaceResultScrape) error {
	args := m.Called(url, model)
	return args.Error(0)
}

// GetPracticeRaceResult is a mock method.
func (m *MockStore) GetPracticeRaceResult(url string, heat, round int) (*models.RaceResultScrape, error) {
	args := m.Called(url, heat, round)

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
func (m *MockStore) GetQualiRaceResult(url string, heat, round int) (*models.RaceResultScrape, error) {
	args := m.Called(url, heat, round)
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
func (m *MockStore) GetFinalRaceResult(url string, final, round int) (*models.RaceResultScrape, error) {
	args := m.Called(url, final, round)
	var res *models.RaceResultScrape
	if args.Get(0) != nil {
		res = args.Get(0).(*models.RaceResultScrape)
	}
	return res, args.Error(1)
}

// MockTracker represents a fake tracker.
type MockTracker struct {
	mock.Mock
}

// EnsureTracking is a mock method.
func (m *MockTracker) EnsureTracking(url string) bool {
	args := m.Called(url)
	return args.Bool(0)
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

	handler := NewHandler(suite.mockStore, suite.mockTracker)

	suite.router = gin.New()
	suite.router.GET("/live", handler.GetLiveTiming)
	suite.router.GET("/practice/round/:round/heat/:heat", handler.GetPracticeRaceResult)
	suite.router.GET("/qualifying/round/:round/heat/:heat", handler.GetQualiRaceResult)
	suite.router.GET("/finals/round/:round/final/:final", handler.GetFinalRaceResult)
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
	req := httptest.NewRequest(http.MethodGet, "/live", nil)
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
	suite.mockTracker.On("EnsureTracking", targetURL).Return(true).Once()
	suite.mockStore.On("GetLiveTiming", targetURL).Return(nil, errors.New("not found")).Once()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/live?target_url="+targetURL, nil)
	suite.router.ServeHTTP(w, req)

	suite.Require().Equal(http.StatusAccepted, w.Code)
	suite.JSONEq(`{
		"status": "success",
		"message": "Starting tracking event, please poll again in a few seconds"
	}`, w.Body.String())
}

func (suite *HandlerTestSuite) TestLive_Success() {
	className := "2 Wheel Drive"
	dummyData := &models.RaceResultScrape{ClassName: &className}

	suite.mockTracker.On("EnsureTracking", targetURL).Return(false).Once()
	suite.mockStore.On("GetLiveTiming", targetURL).Return(dummyData, nil).Once()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/live?target_url="+targetURL, nil)
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
// Practice Heats Tests
// =========================

// TestPractice_MissingURL tests that the API returns a StatusBadRequest 400
// if no target_url query param is passed in.
func (suite *HandlerTestSuite) TestPractice_MissingURL() {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/practice/round/1/heat/2", nil)
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
	suite.mockTracker.On("EnsureTracking", targetURL).Return(true).Once()
	suite.mockStore.On("GetPracticeRaceResult", targetURL, 2, 1).Return(nil, errors.New("not found")).Once()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/practice/round/1/heat/2?target_url="+targetURL, nil)
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
	suite.mockTracker.On("EnsureTracking", targetURL).Return(false).Once()
	suite.mockStore.On("GetPracticeRaceResult", targetURL, 2, 1).Return(nil, errors.New("db error")).Once()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/practice/round/1/heat/2?target_url="+targetURL, nil)
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
	suite.mockTracker.On("EnsureTracking", targetURL).Return(false).Once()
	suite.mockStore.On("GetPracticeRaceResult", targetURL, 2, 1).Return(&models.RaceResultScrape{ClassName: &className}, nil).Once()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/practice/round/1/heat/2?target_url="+targetURL, nil)
	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)
	suite.Contains(w.Body.String(), "LMP2")
}

// =========================
// Qualifying Heats Tests
// =========================

func (suite *HandlerTestSuite) TestQuali_MissingURL() {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/qualifying/round/3/heat/4", nil)
	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusBadRequest, w.Code)
	suite.JSONEq(`{
		"status": "error",
		"message": "Missing required query parameter: target_url"
	}`, w.Body.String())
}

func (suite *HandlerTestSuite) TestQuali_TrackerInitializing() {
	suite.mockTracker.On("EnsureTracking", targetURL).Return(true).Once()
	suite.mockStore.On("GetQualiRaceResult", targetURL, 4, 3).Return(nil, errors.New("not ready")).Once()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/qualifying/round/3/heat/4?target_url="+targetURL, nil)
	suite.router.ServeHTTP(w, req)

	suite.Require().Equal(http.StatusAccepted, w.Code)
	suite.JSONEq(`{
		"status": "success",
		"message": "Starting tracking event, please poll again in a few seconds"
	}`, w.Body.String())
}

func (suite *HandlerTestSuite) TestQuali_StoreError() {
	suite.mockTracker.On("EnsureTracking", targetURL).Return(false).Once()
	suite.mockStore.On("GetQualiRaceResult", targetURL, 4, 3).Return(nil, errors.New("timeout")).Once()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/qualifying/round/3/heat/4?target_url="+targetURL, nil)
	suite.router.ServeHTTP(w, req)

	suite.Require().Equal(http.StatusInternalServerError, w.Code)
	suite.JSONEq(`{
		"status": "error",
		"message": "timeout"
	}`, w.Body.String())
}

func (suite *HandlerTestSuite) TestQuali_Success() {
	className := "GT4"
	suite.mockTracker.On("EnsureTracking", targetURL).Return(false).Once()
	suite.mockStore.On("GetQualiRaceResult", targetURL, 4, 3).Return(&models.RaceResultScrape{ClassName: &className}, nil).Once()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/qualifying/round/3/heat/4?target_url="+targetURL, nil)
	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)
	suite.Contains(w.Body.String(), "GT4")
}

// =========================
// Final Races Tests
// =========================

func (suite *HandlerTestSuite) TestFinal_MissingURL() {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/finals/round/5/final/1", nil)
	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusBadRequest, w.Code)
	suite.JSONEq(`{
		"status": "error",
		"message": "Missing required query parameter: target_url"
	}`, w.Body.String())
}

func (suite *HandlerTestSuite) TestFinal_TrackerInitializing() {
	suite.mockTracker.On("EnsureTracking", targetURL).Return(true).Once()
	suite.mockStore.On("GetFinalRaceResult", targetURL, 1, 5).Return(nil, errors.New("not ready")).Once()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/finals/round/5/final/1?target_url="+targetURL, nil)
	suite.router.ServeHTTP(w, req)

	suite.Require().Equal(http.StatusAccepted, w.Code)
	suite.JSONEq(`{
		"status": "success",
		"message": "Starting tracking event, please poll again in a few seconds"
	}`, w.Body.String())
}

func (suite *HandlerTestSuite) TestFinal_StoreError() {
	suite.mockTracker.On("EnsureTracking", targetURL).Return(false).Once()
	suite.mockStore.On("GetFinalRaceResult", targetURL, 1, 5).Return(nil, errors.New("db disconnect")).Once()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/finals/round/5/final/1?target_url="+targetURL, nil)
	suite.router.ServeHTTP(w, req)

	suite.Require().Equal(http.StatusInternalServerError, w.Code)
	suite.JSONEq(`{
		"status": "error",
		"message": "db disconnect"
	}`, w.Body.String())
}

func (suite *HandlerTestSuite) TestFinal_Success() {
	className := "Hypercar"
	suite.mockTracker.On("EnsureTracking", targetURL).Return(false).Once()
	suite.mockStore.On("GetFinalRaceResult", targetURL, 1, 5).Return(&models.RaceResultScrape{ClassName: &className}, nil).Once()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/finals/round/5/final/1?target_url="+targetURL, nil)
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
