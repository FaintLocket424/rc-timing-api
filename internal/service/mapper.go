package service

import (
	"github.com/FaintLocket424/rc-timing-api/internal/models"
)

// MapCachedHeatResultToDTO converts an internal cached result into its DTO equivalent.
func MapCachedHeatResultToDTO(raw models.CachedHeatResult) models.HeatResultDTO {
	results := make([]models.DriverResultDTO, len(raw.Results))
	for i, res := range raw.Results {
		results[i] = MapCachedDriverResultToDTO(res)
	}

	return models.HeatResultDTO{
		HeatNumber:    raw.HeatNumber,
		RoundNumber:   raw.RoundNumber,
		Class:         raw.Class,
		Results:       results,
		BestLap:       raw.BestLap,
		ClassBestLap:  raw.ClassBestLap,
		ClassBestTime: MapCachedRaceTimeToDTO(raw.ClassBestTime),
	}
}

func MapCachedDriverResultToDTO(raw models.CachedDriverResult) models.DriverResultDTO {
	return models.DriverResultDTO{
		Position:   raw.Position,
		CarNumber:  raw.CarNumber,
		DriverName: raw.DriverName,
		RaceTime:   MapCachedRaceTimeToDTO(raw.RaceTime),
		BestLap:    raw.BestLap,
	}
}

func MapCachedRaceTimeToDTO(raw models.CachedRaceTime) models.RaceTimeDTO {
	return models.RaceTimeDTO{
		Laps:     raw.Laps,
		RaceTime: raw.RaceTime,
	}
}
