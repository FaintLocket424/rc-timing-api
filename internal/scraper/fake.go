package scraper

import (
	"fmt"
	"time"

	"github.com/FaintLocket424/rc-timing-api/internal/models"
)

type FakeScraper struct{}

func (s *FakeScraper) ScrapeResult(baseURL string, heat, round int) (models.RawResult, error) {
	filename := fmt.Sprintf("h%dr%dres.htm", heat, round)
	fmt.Printf("[FakeScraper] Fetching %s/%s...\n", baseURL, filename)

	time.Sleep(500 * time.Millisecond)

	// Mocking some race results
	return models.RawResult{
		HeatNumber:  heat,
		RoundNumber: round,
		ScrapedAt:   time.Now(),
		RawRows: []models.BBKRow{
			{Position: 1, DriverName: "Matthew Peters", Laps: 20, TotalTime: "05:01.2"},
			{Position: 2, DriverName: "Maks Nowak", Laps: 19, TotalTime: "05:04.5"},
			{Position: 3, DriverName: "Oscar Ryley", Laps: 18, TotalTime: "05:03.1"},
		},
	}, nil
}
