package scraper

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/FaintLocket424/rc-timing-api/internal/models"
)

type FakeScraper struct{}

func (s *FakeScraper) ScrapePracticeResult(baseURL string, heat, round int) (models.CachedHeatResult, error) {
	filename := fmt.Sprintf("p%dr%dres.htm", heat, round)
	fmt.Printf("[FakeScraper] Fetching practice %d, round %d at %s/%s...\n", heat, round, baseURL, filename)

	minWait := 150
	maxWait := 750
	waitTime := rand.Intn(maxWait-minWait+1) + minWait
	time.Sleep(time.Duration(waitTime) * time.Millisecond)

	// Mock heat result
	return models.CachedHeatResult{
		HeatNumber:  heat,
		RoundNumber: round,
		Class:       "2wd",
		Results: []models.CachedDriverResult{
			{
				Position:   1,
				CarNumber:  6,
				DriverName: "Matthew Peters",
				RaceTime: models.CachedRaceTime{
					Laps: 16, RaceTime: 5*time.Minute + s.StrToDuration("0.296s"),
				},
				BestLap: s.StrToDuration("14.724s"),
			},
			{
				Position:   2,
				CarNumber:  1,
				DriverName: "Maks Nowak",
				RaceTime: models.CachedRaceTime{
					Laps:     16,
					RaceTime: 5*time.Minute + s.StrToDuration("2.451s"),
				},
				BestLap: s.StrToDuration("14.593s"),
			},
			{
				Position:   3,
				CarNumber:  4,
				DriverName: "Oscar Ryley",
				RaceTime: models.CachedRaceTime{
					Laps:     15,
					RaceTime: 5*time.Minute + s.StrToDuration("2.981s"),
				},
				BestLap: s.StrToDuration("14.236s"),
			},
		},
		BestLap:       s.StrToDuration("14.236s"),
		ClassBestLap:  s.StrToDuration("14.236s"),
		ClassBestTime: models.CachedRaceTime{Laps: 16, RaceTime: 5*time.Minute + s.StrToDuration("0.296s")},
		ScrapedAt:     time.Now(),
	}, nil
}

func (s *FakeScraper) ScrapeQualifyingResult(baseURL string, heat, round int) (models.CachedHeatResult, error) {
	filename := fmt.Sprintf("h%dr%dres.htm", heat, round)
	fmt.Printf("[FakeScraper] Fetching qualifying %d, round %d at %s/%s...\n", heat, round, baseURL, filename)

	minWait := 150
	maxWait := 750
	waitTime := rand.Intn(maxWait-minWait+1) + minWait
	time.Sleep(time.Duration(waitTime) * time.Millisecond)

	// Mock heat result
	return models.CachedHeatResult{
		HeatNumber:  heat,
		RoundNumber: round,
		Class:       "2wd",
		Results: []models.CachedDriverResult{
			{Position: 1, CarNumber: 6, DriverName: "Matthew Peters", RaceTime: models.CachedRaceTime{Laps: 16, RaceTime: 5*time.Minute + s.StrToDuration("0.296s")}, BestLap: s.StrToDuration("14.724s")},
			{Position: 2, CarNumber: 1, DriverName: "Maks Nowak", RaceTime: models.CachedRaceTime{Laps: 16, RaceTime: 5*time.Minute + s.StrToDuration("2.451s")}, BestLap: s.StrToDuration("14.593s")},
			{Position: 3, CarNumber: 4, DriverName: "Oscar Ryley", RaceTime: models.CachedRaceTime{Laps: 15, RaceTime: 5*time.Minute + s.StrToDuration("2.981s")}, BestLap: s.StrToDuration("14.236s")},
		},
		BestLap:       s.StrToDuration("14.236s"),
		ClassBestLap:  s.StrToDuration("14.236s"),
		ClassBestTime: models.CachedRaceTime{Laps: 16, RaceTime: 5*time.Minute + s.StrToDuration("0.296s")},
		ScrapedAt:     time.Now(),
	}, nil
}

func (s *FakeScraper) StrToDuration(str string) time.Duration {
	d, err := time.ParseDuration(str)

	if err != nil {
		return 0
	}

	return d
}
