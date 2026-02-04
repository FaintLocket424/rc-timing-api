package models

import "time"

type Driver struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	CarClass    string `json:"car_class"`
	Transponder string `json:"transponder,omitempty"`
}

// HeatSummary is a tiny version of a heat for the "Schedule" view
type HeatSummary struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Type        string    `json:"type"` // Practice, Qualifying, Final
	Status      string    `json:"status"`
	ScheduledAt time.Time `json:"scheduled_at"`
}

type RaceResult struct {
	Position   int    `json:"pos"`
	DriverName string `json:"name"`
	Laps       int    `json:"laps"`
	TotalTime  string `json:"time"`
}

type Heat struct {
	HeatSummary
	Drivers []Driver     `json:"drivers,omitempty"`
	Results []RaceResult `json:"results,omitempty"`
}

type EventReport struct {
	SourceURL   string    `json:"url"`
	EventName   string    `json:"name"`
	LastUpdated time.Time `json:"updated"`
	Drivers     []Driver  `json:"drivers"`
	Heats       []Heat    `json:"heats"` // We'll store heats as a flat list for easier filtering
}

func (e *EventReport) IsExpired(ttl time.Duration) bool {
	return time.Since(e.LastUpdated) > ttl
}
