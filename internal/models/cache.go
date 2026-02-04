package models

import "time"

type CachedHeatResult struct {
	HeatNumber    int
	RoundNumber   int
	Class         string
	Results       []CachedDriverResult
	BestLap       time.Duration
	ClassBestLap  time.Duration
	ClassBestTime CachedRaceTime
	ScrapedAt     time.Time
}

type CachedDriverResult struct {
	Position   int
	CarNumber  int
	DriverName string
	RaceTime   CachedRaceTime
	BestLap    time.Duration
}

type CachedRaceTime struct {
	Laps     int
	RaceTime time.Duration
}
