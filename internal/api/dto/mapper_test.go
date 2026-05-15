package dto

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/FaintLocket424/opengrid-bridge/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func ptr[T any](v T) *T { return &v }

func TestToRaceResultDTO_FullPracticeResultScrape(t *testing.T) {
	model := &models.RaceResultScrape{
		PracticeNumber: ptr(3),
		Round:          ptr(2),
		ClassName:      ptr("4 Wheel Drive"),
		ElapsedTime:    ptr(2*time.Minute + 12_345*time.Millisecond),
		Drivers: []models.DriverRaceResult{
			{
				CarNumber:       ptr(2),
				Name:            ptr("Dipper Pines"),
				Time:            ptr(2*time.Minute + 10_162*time.Millisecond),
				BestLapDuration: ptr(19_391 * time.Millisecond),
			},
			{
				CarNumber:       ptr(1),
				Name:            ptr("Mabel Pines"),
				Time:            ptr(2*time.Minute + 11_432*time.Millisecond),
				BestLapDuration: ptr(18_998 * time.Millisecond),
			},
		},
		BestLap: &models.RaceBestLap{
			DriverName: ptr("Soos Ramirez"),
			Time:       ptr(18_501 * time.Millisecond),
			LapNumber:  ptr(2),
		},
		ClassFT: &models.ClassFT{
			DriverName: ptr("Grunkle Stan"),
			Laps:       ptr(15),
			Time:       ptr(5*time.Minute + 2_912*time.Millisecond),
			Round:      ptr(2),
		},
		ClassBestLap: &models.ClassBestLap{
			DriverName: ptr("Soos Ramirez"),
			Time:       ptr(18_501 * time.Millisecond),
		},
	}

	dtoResult := ToRaceResultDTO(model)

	bytes, err := json.Marshal(dtoResult)
	require.NoError(t, err)

	expectedJSON := `{
		"practice_number": 3,
		"round": 2,
		"class_name": "4 Wheel Drive",
		"elapsed_ms": 132345,
		"drivers": [
			{
				"car_number": 2,
				"name": "Dipper Pines",
				"time_ms": 130162,
				"best_lap_ms": 19391
			},
			{
				"car_number": 1,
				"name": "Mabel Pines",
				"time_ms": 131432,
				"best_lap_ms": 18998
			}
		],
		"best_lap": {
			"driver_name": "Soos Ramirez",
			"time_ms": 18501,
			"lap_number": 2
		},
		"class_ft": {
			"driver_name": "Grunkle Stan",
			"laps": 15,
			"time_ms": 302912,
			"round": 2
		},
		"class_best_lap": {
			"driver_name": "Soos Ramirez",
			"time_ms": 18501
		}
	}`
	assert.JSONEq(t, expectedJSON, string(bytes), "API Contract has been broken!")
}

func TestToRaceResultDTO_PartialPracticeResultScrape(t *testing.T) {
	model := &models.RaceResultScrape{
		PracticeNumber: ptr(3),
		Round:          ptr(2),
		ClassName:      ptr("4 Wheel Drive"),
		ElapsedTime:    ptr(2*time.Minute + 12_345*time.Millisecond),
		Drivers: []models.DriverRaceResult{
			{
				CarNumber: ptr(2),
				Name:      ptr("Dipper Pines"),
				Time:      ptr(2*time.Minute + 10_162*time.Millisecond),
			},
			{
				CarNumber: ptr(1),
				Name:      ptr("Mabel Pines"),
				Time:      ptr(2*time.Minute + 11_432*time.Millisecond),
			},
		},
		BestLap: &models.RaceBestLap{
			DriverName: ptr("Soos Ramirez"),
			Time:       ptr(18_501 * time.Millisecond),
		},
	}

	dtoResult := ToRaceResultDTO(model)

	bytes, err := json.Marshal(dtoResult)
	require.NoError(t, err)

	expectedJSON := `{
		"practice_number": 3,
		"round": 2,
		"class_name": "4 Wheel Drive",
		"elapsed_ms": 132345,
		"drivers": [
			{
				"car_number": 2,
				"name": "Dipper Pines",
				"time_ms": 130162
			},
			{
				"car_number": 1,
				"name": "Mabel Pines",
				"time_ms": 131432
			}
		],
		"best_lap": {
			"driver_name": "Soos Ramirez",
			"time_ms": 18501
		}
	}`
	assert.JSONEq(t, expectedJSON, string(bytes), "API Contract has been broken!")
}

func TestToRaceResultDTO_FullQualifyingResultScrape(t *testing.T) {
	model := &models.RaceResultScrape{
		HeatNumber:  ptr(3),
		Round:       ptr(2),
		ClassName:   ptr("4 Wheel Drive"),
		ElapsedTime: ptr(2*time.Minute + 12_345*time.Millisecond),
		Drivers: []models.DriverRaceResult{
			{
				CarNumber:       ptr(2),
				Name:            ptr("Dipper Pines"),
				Time:            ptr(2*time.Minute + 10_162*time.Millisecond),
				BestLapDuration: ptr(19_391 * time.Millisecond),
			},
			{
				CarNumber:       ptr(1),
				Name:            ptr("Mabel Pines"),
				Time:            ptr(2*time.Minute + 11_432*time.Millisecond),
				BestLapDuration: ptr(18_998 * time.Millisecond),
			},
		},
		BestLap: &models.RaceBestLap{
			DriverName: ptr("Soos Ramirez"),
			Time:       ptr(18_501 * time.Millisecond),
			LapNumber:  ptr(2),
		},
		ClassFT: &models.ClassFT{
			DriverName: ptr("Grunkle Stan"),
			Laps:       ptr(15),
			Time:       ptr(5*time.Minute + 2_912*time.Millisecond),
			Round:      ptr(2),
		},
		ClassBestLap: &models.ClassBestLap{
			DriverName: ptr("Soos Ramirez"),
			Time:       ptr(18_501 * time.Millisecond),
		},
	}

	dtoResult := ToRaceResultDTO(model)

	bytes, err := json.Marshal(dtoResult)
	require.NoError(t, err)

	expectedJSON := `{
		"heat_number": 3,
		"round": 2,
		"class_name": "4 Wheel Drive",
		"elapsed_ms": 132345,
		"drivers": [
			{
				"car_number": 2,
				"name": "Dipper Pines",
				"time_ms": 130162,
				"best_lap_ms": 19391
			},
			{
				"car_number": 1,
				"name": "Mabel Pines",
				"time_ms": 131432,
				"best_lap_ms": 18998
			}
		],
		"best_lap": {
			"driver_name": "Soos Ramirez",
			"time_ms": 18501,
			"lap_number": 2
		},
		"class_ft": {
			"driver_name": "Grunkle Stan",
			"laps": 15,
			"time_ms": 302912,
			"round": 2
		},
		"class_best_lap": {
			"driver_name": "Soos Ramirez",
			"time_ms": 18501
		}
	}`
	assert.JSONEq(t, expectedJSON, string(bytes), "API Contract has been broken!")
}

func TestToRaceResultDTO_PartialQualifyingResultScrape(t *testing.T) {
	model := &models.RaceResultScrape{
		HeatNumber:  ptr(3),
		Round:       ptr(2),
		ClassName:   ptr("4 Wheel Drive"),
		ElapsedTime: ptr(2*time.Minute + 12_345*time.Millisecond),
		Drivers: []models.DriverRaceResult{
			{
				CarNumber: ptr(2),
				Name:      ptr("Dipper Pines"),
				Time:      ptr(2*time.Minute + 10_162*time.Millisecond),
			},
			{
				CarNumber: ptr(1),
				Name:      ptr("Mabel Pines"),
				Time:      ptr(2*time.Minute + 11_432*time.Millisecond),
			},
		},
		BestLap: &models.RaceBestLap{
			DriverName: ptr("Soos Ramirez"),
			Time:       ptr(18_501 * time.Millisecond),
		},
	}

	dtoResult := ToRaceResultDTO(model)

	bytes, err := json.Marshal(dtoResult)
	require.NoError(t, err)

	expectedJSON := `{
		"heat_number": 3,
		"round": 2,
		"class_name": "4 Wheel Drive",
		"elapsed_ms": 132345,
		"drivers": [
			{
				"car_number": 2,
				"name": "Dipper Pines",
				"time_ms": 130162
			},
			{
				"car_number": 1,
				"name": "Mabel Pines",
				"time_ms": 131432
			}
		],
		"best_lap": {
			"driver_name": "Soos Ramirez",
			"time_ms": 18501
		}
	}`
	assert.JSONEq(t, expectedJSON, string(bytes), "API Contract has been broken!")
}

func TestToRaceResultDTO_FullFinalResultScrape(t *testing.T) {
	model := &models.RaceResultScrape{
		FinalNumber: ptr(3),
		Round:       ptr(2),
		ClassName:   ptr("4 Wheel Drive"),
		ElapsedTime: ptr(2*time.Minute + 12_345*time.Millisecond),
		Drivers: []models.DriverRaceResult{
			{
				CarNumber:       ptr(2),
				Name:            ptr("Dipper Pines"),
				Time:            ptr(2*time.Minute + 10_162*time.Millisecond),
				BestLapDuration: ptr(19_391 * time.Millisecond),
			},
			{
				CarNumber:       ptr(1),
				Name:            ptr("Mabel Pines"),
				Time:            ptr(2*time.Minute + 11_432*time.Millisecond),
				BestLapDuration: ptr(18_998 * time.Millisecond),
			},
		},
		BestLap: &models.RaceBestLap{
			DriverName: ptr("Soos Ramirez"),
			Time:       ptr(18_501 * time.Millisecond),
			LapNumber:  ptr(2),
		},
		ClassFT: &models.ClassFT{
			DriverName: ptr("Grunkle Stan"),
			Laps:       ptr(15),
			Time:       ptr(5*time.Minute + 2_912*time.Millisecond),
			Round:      ptr(2),
		},
		ClassBestLap: &models.ClassBestLap{
			DriverName: ptr("Soos Ramirez"),
			Time:       ptr(18_501 * time.Millisecond),
		},
	}

	dtoResult := ToRaceResultDTO(model)

	bytes, err := json.Marshal(dtoResult)
	require.NoError(t, err)

	expectedJSON := `{
		"final_number": 3,
		"round": 2,
		"class_name": "4 Wheel Drive",
		"elapsed_ms": 132345,
		"drivers": [
			{
				"car_number": 2,
				"name": "Dipper Pines",
				"time_ms": 130162,
				"best_lap_ms": 19391
			},
			{
				"car_number": 1,
				"name": "Mabel Pines",
				"time_ms": 131432,
				"best_lap_ms": 18998
			}
		],
		"best_lap": {
			"driver_name": "Soos Ramirez",
			"time_ms": 18501,
			"lap_number": 2
		},
		"class_ft": {
			"driver_name": "Grunkle Stan",
			"laps": 15,
			"time_ms": 302912,
			"round": 2
		},
		"class_best_lap": {
			"driver_name": "Soos Ramirez",
			"time_ms": 18501
		}
	}`
	assert.JSONEq(t, expectedJSON, string(bytes), "API Contract has been broken!")
}

func TestToRaceResultDTO_PartialFinalResultScrape(t *testing.T) {
	model := &models.RaceResultScrape{
		FinalNumber: ptr(3),
		Round:       ptr(2),
		ClassName:   ptr("4 Wheel Drive"),
		ElapsedTime: ptr(2*time.Minute + 12_345*time.Millisecond),
		Drivers: []models.DriverRaceResult{
			{
				CarNumber: ptr(2),
				Name:      ptr("Dipper Pines"),
				Time:      ptr(2*time.Minute + 10_162*time.Millisecond),
			},
			{
				CarNumber: ptr(1),
				Name:      ptr("Mabel Pines"),
				Time:      ptr(2*time.Minute + 11_432*time.Millisecond),
			},
		},
		BestLap: &models.RaceBestLap{
			DriverName: ptr("Soos Ramirez"),
			Time:       ptr(18_501 * time.Millisecond),
		},
	}

	dtoResult := ToRaceResultDTO(model)

	bytes, err := json.Marshal(dtoResult)
	require.NoError(t, err)

	expectedJSON := `{
		"final_number": 3,
		"round": 2,
		"class_name": "4 Wheel Drive",
		"elapsed_ms": 132345,
		"drivers": [
			{
				"car_number": 2,
				"name": "Dipper Pines",
				"time_ms": 130162
			},
			{
				"car_number": 1,
				"name": "Mabel Pines",
				"time_ms": 131432
			}
		],
		"best_lap": {
			"driver_name": "Soos Ramirez",
			"time_ms": 18501
		}
	}`
	assert.JSONEq(t, expectedJSON, string(bytes), "API Contract has been broken!")
}

func TestToRaceResultDTO_NilHandling(t *testing.T) {
	assert.Nil(t, ToRaceResultDTO(nil))
}

func TestMapToHeatRoundDTOs_Sorting(t *testing.T) {
	// Test the sorting logic for heat/round maps
	input := map[models.HeatRound]struct{}{
		{Round: 2, Heat: 1}: {},
		{Round: 1, Heat: 2}: {},
		{Round: 1, Heat: 1}: {},
	}

	result := mapToHeatRoundDTOs(input)

	require.Len(t, result, 3)
	assert.Equal(t, 1, result[0].Round)
	assert.Equal(t, 1, result[0].Heat)

	assert.Equal(t, 1, result[1].Round)
	assert.Equal(t, 2, result[1].Heat)

	assert.Equal(t, 2, result[2].Round)
	assert.Equal(t, 1, result[2].Heat)
}

func TestToRaceResultsIndexDTO_NilHandling(t *testing.T) {
	assert.Nil(t, ToRaceResultsIndexDTO(nil))
}
