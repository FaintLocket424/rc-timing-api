package service

import (
	"testing"
	"time"

	"github.com/FaintLocket424/rc-timing-api/internal/models"
)

func TestMapCachedHeatResultToDTO(t *testing.T) {
	raw := models.CachedHeatResult{
		HeatNumber:  2,
		RoundNumber: 3,
		Class:       "4wd",
		Results: []models.CachedDriverResult{
			{
				Position:   1,
				CarNumber:  3,
				DriverName: "Matthew Peters",
				RaceTime: models.CachedRaceTime{
					Laps:     20,
					RaceTime: 5*time.Minute + 1342*time.Millisecond,
				},
				BestLap: 15324 * time.Millisecond,
			},
			{
				Position:   2,
				CarNumber:  4,
				DriverName: "Oscar Ryley",
				RaceTime: models.CachedRaceTime{
					Laps:     19,
					RaceTime: 5*time.Minute + 0231*time.Millisecond,
				},
				BestLap: 14324 * time.Millisecond,
			},
		},
		BestLap:      14237 * time.Millisecond,
		ClassBestLap: 12001 * time.Millisecond,
		ClassBestTime: models.CachedRaceTime{
			Laps:     21,
			RaceTime: 5*time.Minute + 1765*time.Millisecond,
		},
		ScrapedAt: time.Now().Add(-10 * time.Minute),
	}

	dto := MapCachedHeatResultToDTO(raw)

	if expectedHeatNumber := 2; dto.HeatNumber != expectedHeatNumber {
		t.Errorf("Expected Heat Number %d, got %d", expectedHeatNumber, dto.HeatNumber)
	}

}
