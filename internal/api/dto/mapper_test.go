package dto

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/FaintLocket424/opengrid-bridge/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestToRaceResultDTO_FullPracticeResultScrape(t *testing.T) {
	model := &models.RaceResultScrape{
		PracticeNumber: ptr(3),
		Round:          ptr(2),
		ClassName:      ptr("4 Wheel Drive"),
		RaceStatus:     ptr("In Progress"),
		ElapsedTime:    ptr(2*time.Minute + 12_345*time.Millisecond),
		Drivers: []models.DriverRaceResult{
			{
				CarNumber:       ptr(2),
				Name:            ptr("Dipper Pines"),
				Time:            ptr(2*time.Minute + 10_162*time.Millisecond),
				BestLapDuration: ptr(19_391 * time.Millisecond),
				WarmupLaps:      ptr(2),
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
		"race_status": "In Progress",
		"elapsed_ms": 132345,
		"drivers": [
			{
				"car_number": 2,
				"name": "Dipper Pines",
				"time_ms": 130162,
				"best_lap_ms": 19391,
				"warmup_laps": 2
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
		RaceStatus:  ptr("In Progress"),
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
		"race_status": "In Progress",
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
		RaceStatus:  ptr("In Progress"),
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
		"race_status": "In Progress",
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

func Test_mapToHeatRoundDTOs(t *testing.T) {
	tests := []struct {
		name  string
		input []models.HeatRound
		want  []HeatRoundDTO
	}{
		{
			name:  "nil input returns empty slice",
			input: nil,
			want:  []HeatRoundDTO{},
		},
		{
			name:  "empty slice returns empty slice",
			input: []models.HeatRound{},
			want:  []HeatRoundDTO{},
		},
		{
			name: "successfully maps multiple elements",
			input: []models.HeatRound{
				{Heat: 1, Round: 2},
				{Heat: 3, Round: 4},
			},
			want: []HeatRoundDTO{
				{Heat: 1, Round: 2},
				{Heat: 3, Round: 4},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mapToHeatRoundDTOs(tt.input)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mapToHeatRoundDTOs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToRaceResultsIndexDTO_NilHandling(t *testing.T) {
	assert.Nil(t, ToRaceResultsIndexDTO(nil))
}

func TestToRaceResultsIndexDTO_ScrapeOmissions(t *testing.T) {
	ts := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	model := &models.RaceResultsIndexScrape{
		Title:     ptr("Event Title"),
		Timestamp: &ts,
		Practice: &models.ResultStatus{
			Results: []models.HeatRound{
				{Heat: 1, Round: 1},
			},
			RoundOveralls: []int{1},
			Overall:       true,
		},
		Qualifying: nil, // Should omit entirely
		Finals: &models.ResultStatus{
			Results:       nil, // Should map to empty arrays
			RoundOveralls: nil,
			Overall:       false,
		},
	}

	dtoResult := ToRaceResultsIndexDTO(model)

	bytes, err := json.Marshal(dtoResult)
	require.NoError(t, err)

	expectedJSON := `{
		"title": "Event Title",
		"timestamp": 1672574400000,
		"practice": {
			"results": [{"h": 1, "r": 1}],
			"round_overalls": [1],
			"overall": true
		},
		"finals": {
			"results": [],
			"round_overalls": [],
			"overall": false
		}
	}`
	assert.JSONEq(t, expectedJSON, string(bytes), "API Contract for RaceResultsIndex has been broken!")
}

// =========================
// Race Schedule Tests
// =========================

func TestToRaceScheduleDTO_Success(t *testing.T) {
	ts := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	startTime1 := time.Date(2023, 1, 1, 13, 0, 0, 0, time.UTC)
	startTime2 := time.Date(2023, 1, 1, 13, 30, 0, 0, time.UTC)

	model := &models.RaceScheduleScrape{
		Title:     ptr("Event Schedule"),
		Timestamp: &ts,
		Practice: []models.ScheduledRace{
			{
				HeatRound: &models.HeatRound{Heat: 1, Round: 1},
				ClassName: ptr("2 Wheel Drive"),
				StartTime: &startTime1,
			},
		},
		Qualifying: []models.ScheduledRace{
			{
				HeatRound: &models.HeatRound{Heat: 2, Round: 1},
				ClassName: ptr("4 Wheel Drive"),
				StartTime: &startTime2,
			},
		},
		Finals: []models.ScheduledRace{},
	}

	dtoResult := ToRaceScheduleDTO(model)

	bytes, err := json.Marshal(dtoResult)
	require.NoError(t, err)

	expectedJSON := `{
		"title": "Event Schedule",
		"timestamp": 1672574400000,
		"practice": [
			{
				"heat": 1,
				"round": 1,
				"class_name": "2 Wheel Drive",
				"start_time": "2023-01-01T13:00:00Z"
			}
		],
		"qualifying": [
			{
				"heat": 2,
				"round": 1,
				"class_name": "4 Wheel Drive",
				"start_time": "2023-01-01T13:30:00Z"
			}
		],
		"finals": []
	}`
	assert.JSONEq(t, expectedJSON, string(bytes), "API Contract for RaceSchedule has been broken!")
}

func TestToRaceScheduleDTO_NilHandling(t *testing.T) {
	assert.Nil(t, ToRaceScheduleDTO(nil))
}

func TestToRaceScheduleDTO_EmptyRaces(t *testing.T) {
	model := &models.RaceScheduleScrape{
		Title:      ptr("Empty Schedule"),
		Practice:   nil,
		Qualifying: nil,
		Finals:     nil,
	}

	dtoResult := ToRaceScheduleDTO(model)

	bytes, err := json.Marshal(dtoResult)
	require.NoError(t, err)

	expectedJSON := `{
		"title": "Empty Schedule",
		"practice": [],
		"qualifying": [],
		"finals": []
	}`
	assert.JSONEq(t, expectedJSON, string(bytes), "API Contract for empty RaceSchedule has been broken!")
}
