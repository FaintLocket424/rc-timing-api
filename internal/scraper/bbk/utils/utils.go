package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	// Matches "11/34.234" or "11/1'34.234"
	resultRegex = regexp.MustCompile(`(?P<laps>-?\d+)/(?:(?P<mins>\d+)')?(?P<secs>\d+\.\d+)`)
	// Matches "11.345" or "11.345[9]"
	lapRegex = regexp.MustCompile(`(?P<time>\d+\.\d+)(?:\[(?P<lap>\d+)\])?`)
)

func NamedCapture(re *regexp.Regexp, input string) map[string]string {
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

// ParseRaceResult parses "11/34.234" or "11/1'34.234"
func ParseRaceResult(res string) (int, time.Duration, error) {
	data := NamedCapture(resultRegex, res)
	if data == nil {
		return 0, 0, fmt.Errorf("invalid result format")
	}

	laps, _ := strconv.Atoi(data["laps"])

	durStr := data["secs"] + "s"
	if data["mins"] != "" {
		durStr = data["mins"] + "m" + durStr
	}

	duration, err := time.ParseDuration(durStr)
	return laps, duration, err
}

// ParseGap parses "+3.45" formats
func ParseGap(gap string) (time.Duration, error) {
	if !strings.HasPrefix(gap, "+") {
		return 0, fmt.Errorf("invalid gap format")
	}
	return time.ParseDuration(gap[1:] + "s")
}

// ParseLap parses "11.345" or "11.345[9]" into duration and lap number
func ParseLap(lapStr string) (duration time.Duration, lapNum int, err error) {
	data := NamedCapture(lapRegex, lapStr)
	if data == nil {
		return 0, 0, fmt.Errorf("invalid lap format")
	}

	duration, _ = time.ParseDuration(data["time"] + "s")
	if data["lap"] != "" {
		lapNum, _ = strconv.Atoi(data["lap"])
	}
	return duration, lapNum, nil
}
