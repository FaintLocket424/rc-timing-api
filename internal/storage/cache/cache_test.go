package cache

import (
	"sync"
	"testing"

	"github.com/FaintLocket424/opengrid-bridge/internal/models"
	"github.com/stretchr/testify/suite"
)

var FakeURL = "http://fake.url"

func ptr[T any](v T) *T {
	return &v
}

type CacheTestSuite struct {
	suite.Suite
	cache *Cache
}

func (suite *CacheTestSuite) SetupTest() {
	suite.cache = NewCache()
}

// --- Live Timing Tests ---

func (suite *CacheTestSuite) TestGetLiveTiming_NotFound() {
	model, err := suite.cache.GetLiveTiming(FakeURL)
	suite.Require().ErrorIs(err, errNotFound)
	suite.Require().Nil(model)
}

func (suite *CacheTestSuite) TestSaveAndGetLiveTiming_Success() {
	data := models.RaceResultScrape{
		HeatNumber: ptr(1),
		ClassName:  ptr("2 Wheel Drive"),
	}

	err := suite.cache.SaveLiveTiming(FakeURL, &data)
	suite.Require().NoError(err)

	model, err := suite.cache.GetLiveTiming(FakeURL)
	suite.Require().NoError(err)
	suite.Require().NotNil(model)
	suite.Require().Equal(&data, model)
}

// --- Practice Results Tests ---

func (suite *CacheTestSuite) TestGetPracticeRaceResult_NotFound() {
	hr := models.HeatRound{Heat: 1, Round: 2}
	model, err := suite.cache.GetPracticeRaceResult(FakeURL, hr)
	suite.Require().ErrorIs(err, errNotFound)
	suite.Require().Nil(model)
}

func (suite *CacheTestSuite) TestSaveAndGetPracticeRaceResult_Success() {
	data := models.RaceResultScrape{
		PracticeNumber: ptr(2),
		Round:          ptr(3),
		ClassName:      ptr("4 Wheel Drive"),
	}

	err := suite.cache.SavePracticeRaceResult(FakeURL, &data)
	suite.Require().NoError(err)

	hr := models.HeatRound{Heat: 2, Round: 3}
	model, err := suite.cache.GetPracticeRaceResult(FakeURL, hr)
	suite.Require().NoError(err)
	suite.Require().NotNil(model)
	suite.Require().Equal(&data, model)
}

func (suite *CacheTestSuite) TestSavePracticeRaceResult_MissingFields() {
	// Missing Round field
	dataNoRound := models.RaceResultScrape{
		PracticeNumber: ptr(2),
	}
	err := suite.cache.SavePracticeRaceResult(FakeURL, &dataNoRound)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "cannot store race result")

	// Missing PracticeNumber field
	dataNoPractice := models.RaceResultScrape{
		Round: ptr(3),
	}
	err = suite.cache.SavePracticeRaceResult(FakeURL, &dataNoPractice)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "cannot store race result")
}

// --- Qualifying Results Tests ---

func (suite *CacheTestSuite) TestGetQualiRaceResult_NotFound() {
	hr := models.HeatRound{Heat: 1, Round: 2}
	model, err := suite.cache.GetQualiRaceResult(FakeURL, hr)
	suite.Require().ErrorIs(err, errNotFound)
	suite.Require().Nil(model)
}

func (suite *CacheTestSuite) TestSaveAndGetQualiRaceResult_Success() {
	data := models.RaceResultScrape{
		HeatNumber: ptr(1),
		Round:      ptr(2),
		ClassName:  ptr("2 Wheel Drive"),
	}

	err := suite.cache.SaveQualiRaceResult(FakeURL, &data)
	suite.Require().NoError(err)

	hr := models.HeatRound{Heat: 1, Round: 2}
	model, err := suite.cache.GetQualiRaceResult(FakeURL, hr)
	suite.Require().NoError(err)
	suite.Require().NotNil(model)
	suite.Require().Equal(&data, model)
}

func (suite *CacheTestSuite) TestSaveQualiRaceResult_MissingFields() {
	// Missing Round field
	dataNoRound := models.RaceResultScrape{
		HeatNumber: ptr(1),
	}
	err := suite.cache.SaveQualiRaceResult(FakeURL, &dataNoRound)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "cannot store race result")

	// Missing HeatNumber field
	dataNoHeat := models.RaceResultScrape{
		Round: ptr(2),
	}
	err = suite.cache.SaveQualiRaceResult(FakeURL, &dataNoHeat)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "cannot store race result")
}

// --- Final Results Tests ---

func (suite *CacheTestSuite) TestGetFinalRaceResult_NotFound() {
	hr := models.HeatRound{Heat: 1, Round: 2}
	model, err := suite.cache.GetFinalRaceResult(FakeURL, hr)
	suite.Require().ErrorIs(err, errNotFound)
	suite.Require().Nil(model)
}

func (suite *CacheTestSuite) TestSaveAndGetFinalRaceResult_Success() {
	data := models.RaceResultScrape{
		FinalNumber: ptr(3),
		Round:       ptr(1),
		ClassName:   ptr("F-Main"),
	}

	err := suite.cache.SaveFinalRaceResult(FakeURL, &data)
	suite.Require().NoError(err)

	hr := models.HeatRound{Heat: 3, Round: 1}
	model, err := suite.cache.GetFinalRaceResult(FakeURL, hr)
	suite.Require().NoError(err)
	suite.Require().NotNil(model)
	suite.Require().Equal(&data, model)
}

func (suite *CacheTestSuite) TestSaveFinalRaceResult_MissingFields() {
	// Missing Round field
	dataNoRound := models.RaceResultScrape{
		FinalNumber: ptr(3),
	}
	err := suite.cache.SaveFinalRaceResult(FakeURL, &dataNoRound)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "cannot store race result")

	// Missing FinalNumber field
	dataNoFinal := models.RaceResultScrape{
		Round: ptr(1),
	}
	err = suite.cache.SaveFinalRaceResult(FakeURL, &dataNoFinal)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "cannot store race result")
}

// --- Results Index Tests ---

func (suite *CacheTestSuite) TestGetRaceResultsIndex_NotFound() {
	model, err := suite.cache.GetRaceResultsIndex(FakeURL)
	suite.Require().ErrorIs(err, errNotFound)
	suite.Require().Nil(model)
}

func (suite *CacheTestSuite) TestSaveAndGetRaceResultsIndex_Success() {
	data := models.RaceResultsIndexScrape{
		Title: ptr("Week 14 results"),
	}

	err := suite.cache.SaveRaceResultsIndex(FakeURL, &data)
	suite.Require().NoError(err)

	model, err := suite.cache.GetRaceResultsIndex(FakeURL)
	suite.Require().NoError(err)
	suite.Require().NotNil(model)
	suite.Require().Equal(&data, model)
}

// --- Race Schedule Tests ---

func (suite *CacheTestSuite) TestGetRaceSchedule_NotFound() {
	model, err := suite.cache.GetRaceSchedule(FakeURL)
	suite.Require().ErrorIs(err, errNotFound)
	suite.Require().Nil(model)
}

func (suite *CacheTestSuite) TestSaveAndGetRaceSchedule_Success() {
	data := models.RaceScheduleScrape{
		Title: ptr("Event Schedule"),
	}

	err := suite.cache.SaveRaceSchedule(FakeURL, &data)
	suite.Require().NoError(err)

	model, err := suite.cache.GetRaceSchedule(FakeURL)
	suite.Require().NoError(err)
	suite.Require().NotNil(model)
	suite.Require().Equal(&data, model)
}

// --- Concurrency Tests ---

func (suite *CacheTestSuite) TestConcurrentAccess() {
	var wg sync.WaitGroup

	data := models.RaceResultScrape{
		HeatNumber: ptr(1),
		ClassName:  ptr("2 Wheel Drive"),
	}

	numGoroutines := 100
	wg.Add(numGoroutines)

	for i := 1; i <= numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()

			if id%2 == 0 {
				for range 100 {
					err := suite.cache.SaveLiveTiming(FakeURL, &data)
					suite.NoError(err, "SaveLiveTiming failed during concurrent write")
				}
			} else {
				for range 100 {
					_, err := suite.cache.GetLiveTiming(FakeURL)
					if err != nil {
						suite.ErrorIs(err, errNotFound, "GetLiveTiming returned an unexpected error")
					}
				}
			}
		}(i)
	}

	wg.Wait()
}

func (suite *CacheTestSuite) TestConcurrentAccess_Mixed() {
	var wg sync.WaitGroup

	numGoroutines := 100
	wg.Add(numGoroutines)

	for i := 1; i <= numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()

			switch id % 4 {
			case 0: // Practice read and write
				for range 50 {
					pData := models.RaceResultScrape{
						PracticeNumber: ptr(id),
						Round:          ptr(1),
					}
					_ = suite.cache.SavePracticeRaceResult(FakeURL, &pData)
					_, _ = suite.cache.GetPracticeRaceResult(FakeURL, models.HeatRound{Heat: id, Round: 1})
				}
			case 1: // Quali read and write
				for range 50 {
					qData := models.RaceResultScrape{
						HeatNumber: ptr(id),
						Round:      ptr(1),
					}
					_ = suite.cache.SaveQualiRaceResult(FakeURL, &qData)
					_, _ = suite.cache.GetQualiRaceResult(FakeURL, models.HeatRound{Heat: id, Round: 1})
				}
			case 2: // Final read and write
				for range 50 {
					fData := models.RaceResultScrape{
						FinalNumber: ptr(id),
						Round:       ptr(1),
					}
					_ = suite.cache.SaveFinalRaceResult(FakeURL, &fData)
					_, _ = suite.cache.GetFinalRaceResult(FakeURL, models.HeatRound{Heat: id, Round: 1})
				}
			case 3: // Schedule and Index writes
				for range 50 {
					_ = suite.cache.SaveRaceResultsIndex(FakeURL, &models.RaceResultsIndexScrape{Title: ptr("Title")})
					_, _ = suite.cache.GetRaceResultsIndex(FakeURL)
					_ = suite.cache.SaveRaceSchedule(FakeURL, &models.RaceScheduleScrape{Title: ptr("Schedule")})
					_, _ = suite.cache.GetRaceSchedule(FakeURL)
				}
			}
		}(i)
	}

	wg.Wait()
}

func TestCacheSuite(t *testing.T) {
	suite.Run(t, new(CacheTestSuite))
}
