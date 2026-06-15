package bbk

import (
	"bytes"
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/FaintLocket424/opengrid-bridge/internal/models"
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
	testFiles := []string{}

	err := filepath.WalkDir("testdata", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && d.Name() == "liveschedule.htm" {
			testFiles = append(testFiles, path)
		}
		return nil
	})
	require.NoError(t, err)

	for _, htmlPath := range testFiles {
		jsonPath := strings.TrimSuffix(htmlPath, filepath.Ext(htmlPath)) + ".json"

		t.Run(htmlPath, func(t *testing.T) {
			if _, err := os.Stat(filepath.Clean(jsonPath)); os.IsNotExist(err) {
				t.Skipf("No golden file found: %s", jsonPath)
			}

			htmlFile, err := os.Open(filepath.Clean(htmlPath))
			require.NoError(t, err)
			defer func() { _ = htmlFile.Close() }()

			actualData, err := parseRaceScheduleHTML(htmlFile)
			require.NoError(t, err)

			expectedJSON, err := os.ReadFile(filepath.Clean(jsonPath))
			require.NoError(t, err)

			var expectedData models.RaceScheduleScrape
			decoder := json.NewDecoder(bytes.NewReader(expectedJSON))
			decoder.DisallowUnknownFields()
			err = decoder.Decode(&expectedData)
			require.NoError(t, err, "Golden JSON is invalid or contains unknown fields")

			// Normalise scrape timestamp
			if expectedData.Timestamp != nil && actualData.Timestamp != nil {
				assert.Equal(t, expectedData.Timestamp.Hour(), actualData.Timestamp.Hour(), "Hour mismatch")
				assert.Equal(t, expectedData.Timestamp.Minute(), actualData.Timestamp.Minute(), "Minute mismatch")
				actualData.Timestamp = expectedData.Timestamp
			}

			// Normalise slices to empty to prevent false-mismatches with nil pointers
			if expectedData.Practice == nil {
				expectedData.Practice = []models.ScheduledRace{}
			}
			if expectedData.Qualifying == nil {
				expectedData.Qualifying = []models.ScheduledRace{}
			}
			if expectedData.Finals == nil {
				expectedData.Finals = []models.ScheduledRace{}
			}

			// Normalise scheduled race start times
			normaliseStartTimes(t, expectedData.Practice, actualData.Practice, "Practice")
			normaliseStartTimes(t, expectedData.Qualifying, actualData.Qualifying, "Qualifying")
			normaliseStartTimes(t, expectedData.Finals, actualData.Finals, "Finals")

			assert.Equal(t, &expectedData, actualData, "Parsed data does not match manually verified golden file")
		})
	}
}
