package service

import (
	"fmt"

	"github.com/FaintLocket424/rc-timing-api/internal/models"
)

func MapRawToDTO(raw models.RawResult) models.HeatResultDTO {
	// We do logic here, like finding the winner
	winner := ""
	if len(raw.RawRows) > 0 {
		winner = raw.RawRows[0].DriverName
	}

	return models.HeatResultDTO{
		HeatName:   fmt.Sprintf("Heat %d", raw.HeatNumber),
		WinnerName: winner,
		Laps:       raw.RawRows[0].Laps,
	}
}
