package scraper

import (
	"time"

	"github.com/FaintLocket424/rc_scraper/internal/models"
)

func Scrape(url string) (models.EventReport, error) {
	time.Sleep(500 * time.Millisecond)

	return models.EventReport{
		SourceURL:   url,
		EventName:   "Trackside GP",
		LastUpdated: time.Now(),
		Drivers: []models.Driver{
			{ID: "101", Name: "FaintLocket", CarClass: "Modified Buggy"},
		},
		Heats: []models.Heat{
			{
				HeatSummary: models.HeatSummary{
					ID: "h1", Name: "Heat 1", Type: "Practice", Status: "Finished",
				},
			},
			{
				HeatSummary: models.HeatSummary{
					ID: "h2", Name: "Heat 2", Type: "Qualifying", Status: "Live",
				},
				Results: []models.RaceResult{
					{Position: 1, DriverName: "FaintLocket", Laps: 5, TotalTime: "2:01.0"},
				},
			},
		},
	}, nil
}
