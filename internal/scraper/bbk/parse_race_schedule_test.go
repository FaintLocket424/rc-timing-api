package bbk

import (
	"os"
	"testing"

	"codeberg.org/OpenGrid-RC/bridge/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func normaliseStartTimes(t *testing.T, expected, actual []models.ScheduledRace, category string) {
	require.Equal(t, len(expected), len(actual), "Number of scheduled %s races mismatch", category)
	for i := range expected {
		expectedRace := &expected[i]
		actualRace := &actual[i]

		if expectedRace.StartTime != nil && actualRace.StartTime != nil {
			assert.Equal(t, expectedRace.StartTime.Hour(), actualRace.StartTime.Hour(), "%s Race %d start hour mismatch", category, i)
			assert.Equal(t, expectedRace.StartTime.Minute(), actualRace.StartTime.Minute(), "%s Race %d start minute mismatch", category, i)
			actualRace.StartTime = expectedRace.StartTime
		}
	}
}

func TestParseRaceSchedule_Golden(t *testing.T) {
	files := collectAndSortTestFiles(t, func(d os.DirEntry) bool {
		return d.Name() == "liveschedule.htm"
	})

	normalise := func(t *testing.T, expected, actual *models.RaceScheduleScrape) {
		if expected.Timestamp != nil && actual.Timestamp != nil {
			assert.Equal(t, expected.Timestamp.Hour(), actual.Timestamp.Hour(), "Hour mismatch")
			assert.Equal(t, expected.Timestamp.Minute(), actual.Timestamp.Minute(), "Minute mismatch")
			actual.Timestamp = expected.Timestamp
		}

		// Normalise slices to empty to prevent false-mismatches with nil pointers
		if expected.Practice == nil {
			expected.Practice = []models.ScheduledRace{}
		}
		if expected.Qualifying == nil {
			expected.Qualifying = []models.ScheduledRace{}
		}
		if expected.Finals == nil {
			expected.Finals = []models.ScheduledRace{}
		}

		// Normalise scheduled race start times
		normaliseStartTimes(t, expected.Practice, actual.Practice, "Practice")
		normaliseStartTimes(t, expected.Qualifying, actual.Qualifying, "Qualifying")
		normaliseStartTimes(t, expected.Finals, actual.Finals, "Finals")
	}

	for _, htmlPath := range files {
		runGoldenTest(t, htmlPath, parseRaceScheduleHTML, normalise)
	}
}
