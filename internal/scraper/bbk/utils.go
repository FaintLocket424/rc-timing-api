package bbk

import (
	"cmp"
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/FaintLocket424/opengrid-bridge/internal/models"
)

// ptr returns a pointer to the value input.
func ptr[T any](v T) *T {
	return &v
}

// atoiPtr parses a string into an int pointer. Returns nil if empty or parsing fails.
func atoiPtr(s string) *int {
	trimmed := strings.TrimSpace(s)
	if trimmed == "" {
		return nil
	}
	v, err := strconv.Atoi(trimmed)
	if err != nil {
		return nil
	}
	return &v
}

// parseTimeToday converts a "15:04" time string to a fully formed *time.Time with today's date.
func parseTimeToday(timeStr string) (*time.Time, error) {
	t, err := time.Parse("15:04", strings.TrimSpace(timeStr))
	if err != nil {
		return nil, err
	}
	now := time.Now()
	date := time.Date(now.Year(), now.Month(), now.Day(), t.Hour(), t.Minute(), 0, 0, now.Location())
	return &date, nil
}

// newResultStatus initializes a models.ResultStatus structure with empty slices.
func newResultStatus() *models.ResultStatus {
	return &models.ResultStatus{
		Results:       []models.HeatRound{},
		RoundOveralls: []int{},
	}
}

// sortHeatRounds sorts a slice of HeatRound structs by Round ascending, then Heat ascending.
func sortHeatRounds(rounds []models.HeatRound) {
	slices.SortFunc(rounds, func(a, b models.HeatRound) int {
		if a.Round != b.Round {
			return cmp.Compare(a.Round, b.Round)
		}
		return cmp.Compare(a.Heat, b.Heat)
	})
}

var (
	// Matches "11/34.234" or "11/1'34.234" patterns.
	resultRegex = regexp.MustCompile(`(?P<laps>-?\d+)/(?:(?P<mins>\d+)')?(?P<secs>\d+\.\d+)`)
	// Matches "11.345" or "11.345[9]" patterns.
	lapRegex = regexp.MustCompile(`(?P<time>\d+\.\d+)(?:\[(?P<lap>\d+)\])?`)
)

// namedCapture extracts named groups from regular expressions and returns them as a map.
func namedCapture(re *regexp.Regexp, input string) map[string]string {
	matches := re.FindStringSubmatch(input)
	if matches == nil {
		return nil
	}
	names := re.SubexpNames()
	results := make(map[string]string)
	for i, name := range names {
		if i != 0 && name != "" {
			results[name] = strings.TrimSpace(matches[i])
		}
	}
	return results
}

// parseRaceResultStr parses "11/34.234" or "11/1'34.234" inputs.
func parseRaceResultStr(res string) (*int, *time.Duration, error) {
	data := namedCapture(resultRegex, res)
	if data == nil {
		return nil, nil, fmt.Errorf("invalid result format")
	}

	laps, err := strconv.Atoi(data["laps"])
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse laps")
	}

	durStr := data["secs"] + "s"
	if data["mins"] != "" {
		durStr = data["mins"] + "m" + durStr
	}

	duration, err := time.ParseDuration(durStr)
	if err != nil {
		return &laps, nil, fmt.Errorf("failed to parse duration")
	}
	return &laps, &duration, nil
}

// parseGapStr parses "+3.45" formats.
func parseGapStr(gap string) (*time.Duration, error) {
	if !strings.HasPrefix(gap, "+") {
		return nil, fmt.Errorf("invalid gap format")
	}
	dur, err := time.ParseDuration(gap[1:] + "s")
	if err != nil {
		return nil, fmt.Errorf("cannot parse gap")
	}

	return &dur, nil
}

// parseLapStr parses "11.345" or "11.345[9]" into duration and lap number.
func parseLapStr(lapStr string) (*time.Duration, *int, error) {
	data := namedCapture(lapRegex, lapStr)
	if data == nil {
		return nil, nil, fmt.Errorf("invalid lap format")
	}

	duration, err := time.ParseDuration(data["time"] + "s")
	if err != nil {
		return nil, nil, fmt.Errorf("cannot parse lap duration")
	}

	if data["lap"] != "" {
		lapNum, err := strconv.Atoi(data["lap"])
		if err != nil {
			return &duration, nil, fmt.Errorf("cannot parse lap")
		}

		return &duration, &lapNum, nil
	}

	return &duration, nil, nil
}
