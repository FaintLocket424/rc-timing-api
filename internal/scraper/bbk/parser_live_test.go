package bbk

func (suite *BBKTestSuite) TestGetLiveTiming_Success() {
	scraper := &BBKScraper{
		Target: suite.server.URL + "/example",
		Client: suite.server.Client(),
	}

	data, err := scraper.GetLiveTiming()

	suite.NoError(err)
	suite.NotNil(data)
	// suite.Equal("John Doe", data.Drivers[0].Name)
}
