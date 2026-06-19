package models

import "time"

// DriverRaceResult represents an individual driver's current or final result
// in a race.
type DriverRaceResult struct {
	CarNumber       *int
	Name            *string
	Laps            *int
	Time            *time.Duration
	BestLapDuration *time.Duration
	BestLapNumber   *int
	LastLapDuration *time.Duration
	AvgLapDuration  *time.Duration
	WarmupLaps      *int
	DNS             *bool
}

// RaceBestLap represents the best lap of a race so far.
type RaceBestLap struct {
	DriverName *string
	Time       *time.Duration
	LapNumber  *int
}

// ClassFT represents the fastest result in the class so far.
type ClassFT struct {
	DriverName *string
	Laps       *int
	Time       *time.Duration
	Round      *int
}

// ClassBestLap represents the fastest lap in the class so far.
type ClassBestLap struct {
	DriverName *string
	Time       *time.Duration
}

// RaceResultScrape represents a full scrape of a race result.
type RaceResultScrape struct {
	PracticeNumber *int
	HeatNumber     *int
	FinalNumber    *int
	Round          *int
	ClassName      *string
	RaceStatus     *string
	ElapsedTime    *time.Duration
	RemainingTime  *time.Duration
	FinishTime     *time.Time
	Drivers        []DriverRaceResult
	BestLap        *RaceBestLap
	ClassFT        *ClassFT
	ClassBestLap   *ClassBestLap
}

// HeatRound represents a specific heat in a specific round.
type HeatRound struct {
	Heat  int
	Round int
}

// ResultStatus represents the races, round overalls and total overalls.
type ResultStatus struct {
	Results       []HeatRound
	RoundOveralls []int
	Overall       bool
}

// RaceResultsIndexScrape represents a full scrape of the results index,
// containing sets of all the practice, quali and finals that have been run.
type RaceResultsIndexScrape struct {
	Title      *string
	Timestamp  *time.Time
	Practice   *ResultStatus
	Qualifying *ResultStatus
	Finals     *ResultStatus
}

// ScheduledRace represents a single race that is scheduled to be run.
type ScheduledRace struct {
	HeatRound *HeatRound
	ClassName *string
	StartTime *time.Time
}

// RaceScheduleScrape represents a full scrape of the race schedule,
// containing all future scheduled races and when they will be run.
type RaceScheduleScrape struct {
	Title      *string
	Timestamp  *time.Time
	Practice   []ScheduledRace
	Qualifying []ScheduledRace
	Finals     []ScheduledRace
}
