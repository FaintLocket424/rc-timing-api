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
	"github.com/PuerkitoBio/goquery"
)

var (
	// Race details regex patterns.
	practiceRegex = regexp.MustCompile(`Practice\s+(?P<practice>\d+)\s+\((?P<class>.*)\)(?:\s+R(?P<round>\d+))?`)
	heatRegex     = regexp.MustCompile(`Heat\s+(?P<heat>\d+)\s+\((?P<class>.*)\)(?:\s+R(?P<round>\d+))?`)
	finalRegex    = regexp.MustCompile(`(?P<final>[A-Z])\.\s+Final\s+\((?P<class>.*)\)(?:\s+L(?P<leg>\d+))?`)

	// Meta and telemetry regex patterns.
	bestLapRegex            = regexp.MustCompile(`Best Lap:\s*(?P<name>.*?)\s+(?P<time>[\d\.]+)\s+L\s*:\s*(?P<lap>\d+)`)
	classFTRegex            = regexp.MustCompile(`Class FT:\s*(?P<name>.*?)\s+(?P<res>[\d\/'\.]+)\s+Av Lap\s+(?P<avg>[\d\.]+)(?:\s+R\s*:\s*(?P<round>\d+))?`)
	classBestLapRegex       = regexp.MustCompile(`Class Best Lap:\s*(?P<name>.*?)\s+(?P<time>\d+\.\d+)`)
	finishedTimeHeaderRegex = regexp.MustCompile(`Finished\s+(?P<finishTime>\d{2}:\d{2})`)
	elapsedTimeHeaderRegex  = regexp.MustCompile(`(?P<elapsed>\d+'\d+s)\s*\((?P<remaining>\d+s)\)`)

	// Index parsing regex patterns.
	resultsLinkRegex       = regexp.MustCompile(`^(?P<type>[phf])(?P<heat>\d+)r(?P<round>\d+)res\.htm$`)
	roundOverallsLinkRegex = regexp.MustCompile(`^(?P<type>[phf])eor(?P<round>\d+)\.htm$`)

	// Utility result/duration parser regex patterns.
	resultRegex = regexp.MustCompile(`(?P<laps>-?\d+)/(?:(?P<mins>\d+)')?(?P<secs>\d+\.\d+)`)
	lapRegex    = regexp.MustCompile(`(?P<time>\d+\.\d+)(?:\[(?P<lap>\d+)\])?`)
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

// findTitleTable searches for the table containing the document's metadata (title and timestamp).
func findTitleTable(doc *goquery.Document) *goquery.Selection {
	var titleTable *goquery.Selection
	doc.Find("table").Each(func(_ int, t *goquery.Selection) {
		if titleTable != nil {
			return
		}
		tds := t.Find("td")
		// Typically Title/Timestamp is in a table with exactly 2 cells where the last is a 24-hour timestamp
		if tds.Length() == 2 {
			timeText := strings.TrimSpace(tds.Last().Text())
			if _, err := parseTimeToday(timeText); err == nil {
				titleTable = t
			}
		}
	})
	return titleTable
}

// parseTitleAndTimestamp extracts Title and Timestamp from the dynamically discovered title table, if it exists.
func parseTitleAndTimestamp(doc *goquery.Document) (*string, *time.Time) {
	titleTable := findTitleTable(doc)
	if titleTable != nil {
		header := titleTable.Find("td")
		if header.Length() > 1 {
			var title *string
			var timestamp *time.Time

			if text := strings.TrimSpace(header.First().Text()); text != "" {
				title = ptr(text)
			}

			if t, err := parseTimeToday(header.Last().Text()); err == nil {
				timestamp = t
			}
			return title, timestamp
		}
	}
	return nil, nil
}

// parseRaceTitleText parses race session type (practice, heat, final), round, and class from a text string.
func parseRaceTitleText(text string) (practice, heat, final, round *int, class *string, matched bool) {
	text = strings.TrimSpace(text)

	if data := namedCapture(practiceRegex, text); data != nil {
		matched = true
		practice = atoiPtr(data["practice"])
		class = ptr(data["class"])
		round = atoiPtr(data["round"])
	} else if data := namedCapture(heatRegex, text); data != nil {
		matched = true
		heat = atoiPtr(data["heat"])
		class = ptr(data["class"])
		round = atoiPtr(data["round"])
	} else if data := namedCapture(finalRegex, text); data != nil {
		matched = true
		finalStr := data["final"]
		if len(finalStr) == 1 {
			r := strings.ToUpper(finalStr)[0]
			final = ptr(int(r - 'A' + 1))
		}
		class = ptr(data["class"])
		if data["leg"] != "" {
			round = atoiPtr(data["leg"])
		} else {
			round = ptr(1)
		}
	}
	return
}

// findHeaderRow searches table rows to find the row containing a target column identifier.
func findHeaderRow(rows *goquery.Selection, target string) (*goquery.Selection, int, bool) {
	var headerRow *goquery.Selection
	headerIdx := -1
	rows.Each(func(i int, r *goquery.Selection) {
		if headerRow != nil {
			return
		}
		r.Children().Each(func(_ int, c *goquery.Selection) {
			if strings.TrimSpace(c.Text()) == target {
				headerRow = r
				headerIdx = i
			}
		})
	})
	return headerRow, headerIdx, headerRow != nil
}

// mapColumnHeaders maps column indices based on matches for expected headers.
func mapColumnHeaders(headerRow *goquery.Selection, mapping map[string][]string) map[string]int {
	idx := make(map[string]int)
	headerRow.Children().Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		for key, matches := range mapping {
			if slices.Contains(matches, text) {
				idx[key] = i
			}
		}
	})
	return idx
}
