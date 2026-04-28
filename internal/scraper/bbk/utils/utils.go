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
func ParseRaceResult(res string) (*int, *time.Duration, error) {
	data := NamedCapture(resultRegex, res)
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

// ParseGap parses "+3.45" formats
func ParseGap(gap string) (*time.Duration, error) {
	if !strings.HasPrefix(gap, "+") {
		return nil, fmt.Errorf("invalid gap format")
	}
	dur, err := time.ParseDuration(gap[1:] + "s")
	if err != nil {
		return nil, fmt.Errorf("cannot parse gap")
	}

	return &dur, nil
}

// ParseLap parses "11.345" or "11.345[9]" into duration and lap number
func ParseLap(lapStr string) (*time.Duration, *int, error) {
	data := NamedCapture(lapRegex, lapStr)
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
		} else {
			return &duration, &lapNum, nil
		}
	} else {
		return &duration, nil, nil
	}
}
