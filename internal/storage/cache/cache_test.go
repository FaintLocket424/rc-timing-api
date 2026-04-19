package cache

import (
	"sync"
	"testing"

	"github.com/FaintLocket424/rc-timing-api/internal/models"
	"github.com/stretchr/testify/suite"
)

type CacheTestSuite struct {
	suite.Suite
	cache *Cache
}

func (suite *CacheTestSuite) SetupTest() {
	suite.cache = NewCache()
}

func (suite *CacheTestSuite) TestGetLiveTiming_NotFound() {
	model, err := suite.cache.GetLiveTiming("http://fake.url")
	suite.Require().Error(err)
	suite.Require().Nil(model)
}

func (suite *CacheTestSuite) TestSaveAndGetLiveTiming_Success() {
	data := models.LiveTimingScrape{
		HeatNumber: 1,
		ClassName:  "2 Wheel Drive",
	}

	url := "http://fake.url"

	err := suite.cache.SaveLiveTiming(url, &data)
	suite.Require().NoError(err)

	model, err := suite.cache.GetLiveTiming(url)
	suite.Require().NoError(err)
	suite.Require().NotNil(model)
	suite.Require().Equal(&data, model)
}

func (suite *CacheTestSuite) TestConcurrentAccess() {
	var wg sync.WaitGroup

	data := models.LiveTimingScrape{
		HeatNumber: 1,
		ClassName:  "2 Wheel Drive",
	}

	url := "http://fake.url"
	numGoroutines := 100
	wg.Add(numGoroutines)

	for i := 1; i <= numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()

			if id%2 == 0 {
				for range 100 {
					suite.cache.SaveLiveTiming(url, &data)
				}
			} else {
				for range 100 {
					suite.cache.GetLiveTiming(url)
				}
			}
		}(i)
	}

	wg.Wait()

}

func TestCacheSuite(t *testing.T) {
	suite.Run(t, new(CacheTestSuite))
}
