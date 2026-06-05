package dto

// DriverRaceResult represents an individual driver's current or final result
// in a race.
type DriverRaceResult struct {
	CarNumber     *int    `json:"car_number,omitempty"`
	Name          *string `json:"name,omitempty"`
	Laps          *int    `json:"laps,omitempty"`
	TimeMs        *int64  `json:"time_ms,omitempty"`
	BestLapMs     *int64  `json:"best_lap_ms,omitempty"`
	BestLapNumber *int    `json:"best_lap_number,omitempty"`
	LastLapMs     *int64  `json:"last_lap_ms,omitempty"`
	WarmupLaps    *int    `json:"warmup_laps,omitempty"`
}

// RaceBestLap represents the best lap of a race so far.
type RaceBestLap struct {
	DriverName *string `json:"driver_name,omitempty"`
	TimeMs     *int64  `json:"time_ms,omitempty"`
	LapNumber  *int    `json:"lap_number,omitempty"`
}

// ClassFT represents the fastest result in the class so far.
type ClassFT struct {
	DriverName *string `json:"driver_name,omitempty"`
	Laps       *int    `json:"laps,omitempty"`
	TimeMs     *int64  `json:"time_ms,omitempty"`
	Round      *int    `json:"round,omitempty"`
}

// ClassBestLap represents the fastest lap in the class so far.
type ClassBestLap struct {
	DriverName *string `json:"driver_name,omitempty"`
	TimeMs     *int64  `json:"time_ms,omitempty"`
}

// RaceResultScrapeDTO represents a full scrape of a race result.
type RaceResultScrapeDTO struct {
	PracticeNumber *int               `json:"practice_number,omitempty"`
	HeatNumber     *int               `json:"heat_number,omitempty"`
	FinalNumber    *int               `json:"final_number,omitempty"`
	ClassName      *string            `json:"class_name,omitempty"`
	RaceStatus     *string            `json:"race_status,omitempty"`
	Round          *int               `json:"round,omitempty"`
	ElapsedTimeMs  *int64             `json:"elapsed_ms,omitempty"`
	Drivers        []DriverRaceResult `json:"drivers,omitempty"`
	BestLap        *RaceBestLap       `json:"best_lap,omitempty"`
	ClassFT        *ClassFT           `json:"class_ft,omitempty"`
	ClassBestLap   *ClassBestLap      `json:"class_best_lap,omitempty"`
}

// HeatRoundDTO represents a specific heat in a specific round.
type HeatRoundDTO struct {
	Heat  int `json:"h"`
	Round int `json:"r"`
}

// ResultStatusDTO represents the races, round overalls and total overalls.
type ResultStatusDTO struct {
	Results       []HeatRoundDTO `json:"results"`
	RoundOveralls []int          `json:"round_overalls"`
	Overall       bool           `json:"overall"`
}

// RaceResultsIndexScrapeDTO represents a full scrape of the results index,
// containing sets of all the practice, quali and finals that have been run.
type RaceResultsIndexScrapeDTO struct {
	Title      *string          `json:"title,omitempty"`
	Timestamp  *int64           `json:"timestamp,omitempty"`
	Practice   *ResultStatusDTO `json:"practice,omitempty"`
	Qualifying *ResultStatusDTO `json:"qualifying,omitempty"`
	Finals     *ResultStatusDTO `json:"finals,omitempty"`
}
