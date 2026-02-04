package models

import "time"

type RawResult struct {
	HeatNumber  int
	RoundNumber int
	RawRows     []BBKRow // All data, even the "ugly" stuff
	ScrapedAt   time.Time
}

type BBKRow struct {
	Position    int
	DriverName  string
	Transponder string // We keep this internally even if the API doesn't show it
	Laps        int
	LastLapTime string
	TotalTime   string
}
