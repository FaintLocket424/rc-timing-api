package models

import "time"

type DriverRaceResult struct {
	CarNumber       *int
	Name            *string
	Laps            *int
	Time            *time.Duration
	BestLapDuration *time.Duration
	BestLapNumber   *int
	LastLapDuration *time.Duration
}

type BestLap struct {
	DriverName *string
	Time       *time.Duration
	LapNumber  *int
}

type ClassFT struct {
	DriverName     *string
	Laps           *int
	Time           *time.Duration
	AvgLapDuration *time.Duration
	Round          *int
}

type ClassBestLap struct {
	DriverName *string
	Time       *time.Duration
}

type ResultScrape struct {
	HeatNumber   *int
	ClassName    *string
	Round        *int
	ElapsedTime  *time.Duration
	Drivers      []DriverRaceResult
	BestLap      *BestLap
	ClassFT      *ClassFT
	ClassBestLap *ClassBestLap
}
