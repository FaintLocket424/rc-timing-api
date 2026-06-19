package bbk

import (
	"os"
	"regexp"
	"testing"

	"codeberg.org/OpenGrid-RC/bridge/internal/models"
	"github.com/stretchr/testify/assert"
)

// Matches: liveraceres.htm, p1r1res.htm, h5r2res.htm, f12r1res.htm, etc.
var compatibleFileRegex = regexp.MustCompile(`(?i)^(liveraceres|[phf]\d+r\d+res)\.htm$`)

func TestParseRaceResult_Golden(t *testing.T) {
	files := collectAndSortTestFiles(t, func(d os.DirEntry) bool {
		return compatibleFileRegex.MatchString(d.Name())
	})

	normalise := func(t *testing.T, expected, actual *models.RaceResultScrape) {
		if expected.FinishTime != nil && actual.FinishTime != nil {
			assert.Equal(t, expected.FinishTime.Hour(), actual.FinishTime.Hour(), "Hour mismatch")
			assert.Equal(t, expected.FinishTime.Minute(), actual.FinishTime.Minute(), "Minute mismatch")
			actual.FinishTime = expected.FinishTime
		}
	}

	for _, htmlPath := range files {
		runGoldenTest(t, htmlPath, parseRaceResultHTML, normalise)
	}
}
