package bbk

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"codeberg.org/OpenGrid-RC/bridge/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Matches: liveraceres.htm, p1r1res.htm, h5r2res.htm, f12r1res.htm, etc.
var compatibleFileRegex = regexp.MustCompile(`(?i)^(liveraceres|[phf]\d+r\d+res)\.htm$`)

func TestParseRaceResult_Golden(t *testing.T) {
	testFiles := []string{}

	err := filepath.WalkDir("testdata", func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && compatibleFileRegex.MatchString(d.Name()) {
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

			actualData, err := parseRaceResultHTML(htmlFile)
			require.NoError(t, err)

			expectedJSON, err := os.ReadFile(filepath.Clean(jsonPath))
			require.NoError(t, err)

			var expectedData models.RaceResultScrape
			decoder := json.NewDecoder(bytes.NewReader(expectedJSON))
			decoder.DisallowUnknownFields()
			err = decoder.Decode(&expectedData)
			require.NoError(t, err, "Golden JSON is invalid or contains unknown fields")

			// Normalise scrape finish time safely
			if expectedData.FinishTime != nil && actualData.FinishTime != nil {
				assert.Equal(t, expectedData.FinishTime.Hour(), actualData.FinishTime.Hour(), "Hour mismatch")
				assert.Equal(t, expectedData.FinishTime.Minute(), actualData.FinishTime.Minute(), "Minute mismatch")
				actualData.FinishTime = expectedData.FinishTime
			}

			assert.Equal(t, &expectedData, actualData, "Parsed data does not match manually verified golden file")
		})
	}
}
