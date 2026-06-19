package bbk

import (
	"os"
	"testing"

	"codeberg.org/OpenGrid-RC/bridge/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestParseRaceResultsIndex_Golden(t *testing.T) {
	files := collectAndSortTestFiles(t, func(d os.DirEntry) bool {
		return d.Name() == "liveresults.htm"
	})

	normalise := func(t *testing.T, expected, actual *models.RaceResultsIndexScrape) {
		if expected.Timestamp != nil && actual.Timestamp != nil {
			assert.Equal(t, expected.Timestamp.Hour(), actual.Timestamp.Hour(), "Hour mismatch")
			assert.Equal(t, expected.Timestamp.Minute(), actual.Timestamp.Minute(), "Minute mismatch")
			actual.Timestamp = expected.Timestamp
		}
	}

	for _, htmlPath := range files {
		runGoldenTest(t, htmlPath, parseRaceResultsIndexHTML, normalise)
	}
}
