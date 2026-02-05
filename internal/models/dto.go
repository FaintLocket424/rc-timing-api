package models

import "time"

type HeatResultDTO struct {
	HeatNumber    int               `json:"heat_number"`
	RoundNumber   int               `json:"round_number"`
	Class         string            `json:"class"`
	Results       []DriverResultDTO `json:"results"`
	BestLap       time.Duration     `json:"best_lap"`
	ClassBestLap  time.Duration     `json:"class_best_lap"`
	ClassBestTime RaceTimeDTO       `json:"class_best_time"`
}

type DriverResultDTO struct {
	Position   int           `json:"position"`
	CarNumber  int           `json:"car_number"`
	DriverName string        `json:"driver_name"`
	RaceTime   RaceTimeDTO   `json:"race_time"`
	BestLap    time.Duration `json:"best_lap"`
}

type RaceTimeDTO struct {
	Laps     int           `json:"laps"`
	RaceTime time.Duration `json:"race_time"`
}

type EventMetaDTO struct {
	NumCompetitors int `json:"num_competitors"`
}
