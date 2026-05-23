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

func TestParseRaceResultsIndex_Golden(t *testing.T) {
	testFiles := []string{}

	err := filepath.WalkDir("testdata", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && d.Name() == "liveresults.htm" {
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

			actualData, err := parseRaceResultsIndexHTML(htmlFile)
			require.NoError(t, err)

			expectedJSON, err := os.ReadFile(filepath.Clean(jsonPath))
			require.NoError(t, err)

			var expectedData models.RaceResultsIndexScrape
			decoder := json.NewDecoder(bytes.NewReader(expectedJSON))
			decoder.DisallowUnknownFields()
			err = decoder.Decode(&expectedData)
			require.NoError(t, err, "Golden JSON is invalid or contains unknown fields")

			assert.Equal(t, expectedData.Timestamp.Hour(), actualData.Timestamp.Hour(), "Hour mismatch")
			assert.Equal(t, expectedData.Timestamp.Minute(), actualData.Timestamp.Minute(), "Minute mismatch")
			actualData.Timestamp = expectedData.Timestamp

			assert.Equal(t, &expectedData, actualData, "Parsed data does not match manually verified golden file")
		})
	}
}
