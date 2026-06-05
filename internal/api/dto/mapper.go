// Package dto handles mapping internal data structure to "Data Transfer Objects (DTOs)".
// The DTOs form the API interface.
package dto

import (
	"time"

	"github.com/FaintLocket424/opengrid-bridge/internal/models"
)

// durationToMs safely converts a pointer of time.Duration to a pointer of int64 milliseconds.
func durationToMs(d *time.Duration) *int64 {
	if d == nil {
		return nil
	}
	ms := d.Milliseconds()
	return &ms
}

// ToRaceResultDTO converts the internal RaceResultScrape model to the external DTO.
func ToRaceResultDTO(m *models.RaceResultScrape) *RaceResultScrapeDTO {
	if m == nil {
		return nil
	}

	var drivers []DriverRaceResult
	if len(m.Drivers) > 0 {
		drivers = make([]DriverRaceResult, len(m.Drivers))
		for i, d := range m.Drivers {
			drivers[i] = DriverRaceResult{
				CarNumber:     d.CarNumber,
				Name:          d.Name,
				Laps:          d.Laps,
				TimeMs:        durationToMs(d.Time),
				BestLapMs:     durationToMs(d.BestLapDuration),
				BestLapNumber: d.BestLapNumber,
				LastLapMs:     durationToMs(d.LastLapDuration),
				WarmupLaps:    d.WarmupLaps,
			}
		}
	}

	var bestLap *RaceBestLap
	if m.BestLap != nil {
		bestLap = &RaceBestLap{
			DriverName: m.BestLap.DriverName,
			TimeMs:     durationToMs(m.BestLap.Time),
			LapNumber:  m.BestLap.LapNumber,
		}
	}

	var classFT *ClassFT
	if m.ClassFT != nil {
		classFT = &ClassFT{
			DriverName: m.ClassFT.DriverName,
			Laps:       m.ClassFT.Laps,
			TimeMs:     durationToMs(m.ClassFT.Time),
			Round:      m.ClassFT.Round,
		}
	}

	var classBestLap *ClassBestLap
	if m.ClassBestLap != nil {
		classBestLap = &ClassBestLap{
			DriverName: m.ClassBestLap.DriverName,
			TimeMs:     durationToMs(m.ClassBestLap.Time),
		}
	}

	return &RaceResultScrapeDTO{
		PracticeNumber: m.PracticeNumber,
		HeatNumber:     m.HeatNumber,
		FinalNumber:    m.FinalNumber,
		ClassName:      m.ClassName,
		RaceStatus:     m.RaceStatus,
		Round:          m.Round,
		ElapsedTimeMs:  durationToMs(m.ElapsedTime),
		Drivers:        drivers,
		BestLap:        bestLap,
		ClassFT:        classFT,
		ClassBestLap:   classBestLap,
	}
}

// ToRaceResultsIndexDTO converts the internal RaceResultsIndexScrape model to the external DTO.
func ToRaceResultsIndexDTO(m *models.RaceResultsIndexScrape) *RaceResultsIndexScrapeDTO {
	if m == nil {
		return nil
	}

	data := &RaceResultsIndexScrapeDTO{
		Title:      m.Title,
		Practice:   mapToResultStatusDTO(m.Practice),
		Qualifying: mapToResultStatusDTO(m.Qualifying),
		Finals:     mapToResultStatusDTO(m.Finals),
	}

	if m.Timestamp != nil {
		ts := m.Timestamp.UnixMilli()
		data.Timestamp = &ts
	}

	return data
}

// mapToResultStatusDTO is a helper to convert a models.ResultStatus to a ResultStatusDTO.
func mapToResultStatusDTO(r *models.ResultStatus) *ResultStatusDTO {
	if r == nil {
		return nil
	}

	var roundOveralls []int
	if len(r.RoundOveralls) > 0 {
		roundOveralls = r.RoundOveralls
	}

	return &ResultStatusDTO{
		Results:       mapToHeatRoundDTOs(r.Results),
		RoundOveralls: roundOveralls,
		Overall:       r.Overall,
	}
}

// mapToHeatRoundDTOs safely converts the internal struct map to a JSON-friendly array.
func mapToHeatRoundDTOs(m []models.HeatRound) []HeatRoundDTO {
	if len(m) == 0 {
		return nil
	}

	res := make([]HeatRoundDTO, len(m))
	for i, k := range m {
		res[i] = HeatRoundDTO{Heat: k.Heat, Round: k.Round}
	}

	return res
}
