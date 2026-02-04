package models

type HeatResultDTO struct {
	HeatName   string `json:"heat_name"`
	WinnerName string `json:"winner"`
	Laps       int    `json:"laps"`
}
