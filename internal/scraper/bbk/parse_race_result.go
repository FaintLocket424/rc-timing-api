package bbk

import (
	"fmt"
	"io"
	"log/slog"
	"strings"
	"time"

	"codeberg.org/OpenGrid-RC/bridge/internal/models"
	"github.com/PuerkitoBio/goquery"
)

const driverColName = "Driver"

func parseRaceResultHTML(body io.Reader) (*models.RaceResultScrape, error) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML body: %w", err)
	}

	scrape := &models.RaceResultScrape{}

	tables := doc.Find("table")

	var headerTable, driversTable, metaTable *goquery.Selection

	// Dynamically identify tables based on structural signatures
	tables.Each(func(_ int, t *goquery.Selection) {
		hasHeaderPattern := false
		hasDriverHeaders := false
		hasMetaPattern := false

		t.Find("td").Each(func(_ int, td *goquery.Selection) {
			txt := strings.TrimSpace(td.Text())
			if practiceRegex.MatchString(txt) || heatRegex.MatchString(txt) || finalRegex.MatchString(txt) {
				hasHeaderPattern = true
			}
			if txt == driverColName {
				hasDriverHeaders = true
			}
			if strings.HasPrefix(txt, "Best Lap:") || strings.HasPrefix(txt, "Class FT:") || strings.HasPrefix(txt, "Class Best Lap:") || td.HasClass("fastest-lap") {
				hasMetaPattern = true
			}
		})

		switch {
		case hasHeaderPattern:
			headerTable = t
		case hasDriverHeaders:
			driversTable = t
		case hasMetaPattern:
			metaTable = t
		}
	})

	// Safely execute parsing steps if the corresponding structures exist
	if headerTable != nil {
		parseHeader(headerTable, scrape)
	}
	if driversTable != nil {
		parseDrivers(driversTable, scrape)
	}
	if metaTable != nil {
		parseMeta(metaTable, scrape)
	}

	const tolerance = 100 * time.Millisecond

	// 1. Individual Driver Data Integrity Check
	for _, dr := range scrape.Drivers {
		if dr.AvgLapDuration != nil && dr.Time != nil && dr.Laps != nil && *dr.Laps > 0 {
			expectedAvg := *dr.Time / time.Duration(*dr.Laps)
			diff := expectedAvg - *dr.AvgLapDuration
			if diff < 0 {
				diff = -diff
			}

			if diff > tolerance {
				driverName := "unknown"
				if dr.Name != nil {
					driverName = *dr.Name
				}
				return nil, fmt.Errorf("data integrity check failed for driver %q: reported average lap %s, calculated average lap %s (diff %s exceeds tolerance of %s)",
					driverName, dr.AvgLapDuration, expectedAvg, diff, tolerance)
			}
		}
	}

	// 2. Class FT Data Integrity Check
	if scrape.ClassFT != nil && scrape.ClassFT.AvgLapDuration != nil && scrape.ClassFT.Time != nil && scrape.ClassFT.Laps != nil && *scrape.ClassFT.Laps > 0 {
		expectedAvg := *scrape.ClassFT.Time / time.Duration(*scrape.ClassFT.Laps)
		diff := expectedAvg - *scrape.ClassFT.AvgLapDuration
		if diff < 0 {
			diff = -diff
		}

		if diff > tolerance {
			driverName := "unknown"
			if scrape.ClassFT.DriverName != nil {
				driverName = *scrape.ClassFT.DriverName
			}
			return nil, fmt.Errorf("data integrity check failed for Class FT (%s): reported average lap %s, calculated average lap %s (diff %s exceeds tolerance of %s)",
				driverName, scrape.ClassFT.AvgLapDuration, expectedAvg, diff, tolerance)
		}
	}

	// 3. Global Best Lap Data Integrity Check
	if scrape.BestLap != nil && scrape.BestLap.DriverName != nil && scrape.BestLap.Time != nil {
		matched := false
		bestLapDriverName := *scrape.BestLap.DriverName
		bestLapTime := *scrape.BestLap.Time

		for _, dr := range scrape.Drivers {
			if dr.Name != nil && dr.BestLapDuration != nil {
				nameMatches := strings.EqualFold(strings.TrimSpace(*dr.Name), strings.TrimSpace(bestLapDriverName))
				timeMatches := *dr.BestLapDuration == bestLapTime

				lapMatches := true
				if scrape.BestLap.LapNumber != nil && dr.BestLapNumber != nil {
					lapMatches = *scrape.BestLap.LapNumber == *dr.BestLapNumber
				}

				if nameMatches && timeMatches && lapMatches {
					matched = true
					break
				}
			}
		}

		if !matched {
			return nil, fmt.Errorf("data integrity check failed: global best lap by %q (%s) does not match any driver's parsed best lap in the race",
				bestLapDriverName, bestLapTime)
		}
	}

	return scrape, nil
}

func parseHeader(s *goquery.Selection, lt *models.RaceResultScrape) {
	s.Find("td").Each(func(_ int, td *goquery.Selection) {
		text := strings.TrimSpace(td.Text())

		p, h, f, r, c, matched := parseRaceTitleText(text)
		if matched {
			lt.PracticeNumber = p
			lt.HeatNumber = h
			lt.FinalNumber = f
			lt.Round = r
			lt.ClassName = c
		}

		if finishMatch := namedCapture(finishedTimeHeaderRegex, text); finishMatch != nil {
			if t, err := parseTimeToday(finishMatch["finishTime"]); err == nil {
				lt.FinishTime = t
			}
			lt.RaceStatus = ptr("Finished")
		}

		if durMatch := namedCapture(elapsedTimeHeaderRegex, text); durMatch != nil {
			elapsedRaw := strings.Replace(durMatch["elapsed"], "'", "m", 1)
			remainingRaw := durMatch["remaining"]

			if elapsed, err := time.ParseDuration(elapsedRaw); err == nil {
				lt.ElapsedTime = ptr(elapsed)
			}

			if remaining, err := time.ParseDuration(remainingRaw); err == nil {
				lt.RemainingTime = ptr(remaining)
			}
			lt.RaceStatus = ptr("In Progress")
		}
	})
}

func parseDrivers(drivers *goquery.Selection, lt *models.RaceResultScrape) {
	rows := drivers.Find("tr")
	if rows.Length() < 2 {
		slog.Error("drivers table missing data rows")
		return
	}

	headerRow, headerIdx, found := findHeaderRow(rows, driverColName)
	if !found {
		slog.Error("drivers table missing column headers")
		return
	}

	idx := mapColumnHeaders(headerRow, map[string][]string{
		"car":  {"C", "B"},
		"name": {driverColName},
		"res":  {"Result"},
		"best": {"B-Lap", "B-Lp"},
		"last": {"L-Lap", "L-Lp"},
		"avg":  {"A-Lap", "A-Lp"},
	})

	// Process rows
	rows.Slice(headerIdx+1, rows.Length()).Each(func(i int, s *goquery.Selection) {
		cols := s.Children()
		dr := models.DriverRaceResult{}

		// Driver name
		if nameIdx, ok := idx["name"]; ok && nameIdx < cols.Length() {
			dr.Name = ptr(strings.TrimSpace(cols.Eq(nameIdx).Text()))
		}

		// Skip empty spacer/filler rows to keep data clean
		if dr.Name == nil || *dr.Name == "" {
			return
		}

		// Car Number
		if carIdx, ok := idx["car"]; ok && carIdx < cols.Length() {
			if v := atoiPtr(cols.Eq(carIdx).Text()); v != nil {
				dr.CarNumber = v
			}
		}

		// Result
		if resIdx, ok := idx["res"]; ok && resIdx < cols.Length() {
			resTd := cols.Eq(resIdx)
			resText := strings.TrimSpace(resTd.Text())

			// Determine if the driver has achieved a personal/race improvement
			dr.Improvement = ptr(resTd.HasClass("livehtml-racebettertime"))

			if laps, dur, err := parseRaceResultStr(resText); err == nil {
				dr.Laps = laps
				if dr.Laps != nil && *dr.Laps < 0 && len(lt.Drivers) > 0 && lt.Drivers[0].Laps != nil {
					*dr.Laps += *lt.Drivers[0].Laps
				}
				dr.Time = dur
			} else if gap, err := parseGapStr(resText); err == nil {
				if len(lt.Drivers) > 0 && lt.Drivers[0].Laps != nil && lt.Drivers[0].Time != nil {
					dr.Laps = ptr(*lt.Drivers[0].Laps)
					dr.Time = ptr(*lt.Drivers[0].Time + *gap)
				}
			} else if numStr, ok := strings.CutPrefix(resText, "W-"); ok {
				if v := atoiPtr(numStr); v != nil {
					dr.WarmupLaps = v
				}
			} else if resText == "DNS" {
				dr.DNS = ptr(true) // Set DNS to true using the package-level helper
			} else if resText != "" && resText != "DNS" && resText != "DNF" && resText != "DSQ" {
				slog.Warn("unparseable result", "row", i+1, "result", resText)
			}
		}

		// Best Lap duration
		if bestIdx, ok := idx["best"]; ok && bestIdx < cols.Length() {
			text := strings.TrimSpace(cols.Eq(bestIdx).Text())
			if dur, ln, err := parseLapStr(text); err == nil {
				dr.BestLapDuration = dur
				dr.BestLapNumber = ln
			}
		}

		// Average Lap duration
		if avgIdx, ok := idx["avg"]; ok && avgIdx < cols.Length() {
			text := strings.TrimSpace(cols.Eq(avgIdx).Text())
			if dur, _, err := parseLapStr(text); err == nil {
				dr.AvgLapDuration = dur
			}
		}

		// Last Lap duration
		if lastIdx, ok := idx["last"]; ok && lastIdx < cols.Length() {
			text := strings.TrimSpace(cols.Eq(lastIdx).Text())
			if dur, _, err := parseLapStr(text); err == nil {
				dr.LastLapDuration = dur
			}
		}

		lt.Drivers = append(lt.Drivers, dr)
	})
}

func parseMeta(meta *goquery.Selection, lt *models.RaceResultScrape) {
	meta.Find("tr td").Each(func(_ int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())

		switch {
		case strings.HasPrefix(text, "Best Lap:"):
			if data := namedCapture(bestLapRegex, text); data != nil {
				lt.BestLap = &models.RaceBestLap{DriverName: ptr(data["name"])}
				dur, ln, _ := parseLapStr(fmt.Sprintf("%s[%s]", data["time"], data["lap"]))
				lt.BestLap.Time, lt.BestLap.LapNumber = dur, ln
			} else {
				slog.Warn("failed to parse Best Lap meta", "raw", text)
			}
		case strings.HasPrefix(text, "Class FT:"):
			if data := namedCapture(classFTRegex, text); data != nil {
				lt.ClassFT = &models.ClassFT{DriverName: ptr(data["name"])}
				laps, dur, _ := parseRaceResultStr(data["res"])
				lt.ClassFT.Laps, lt.ClassFT.Time = laps, dur
				if r := atoiPtr(data["round"]); r != nil {
					lt.ClassFT.Round = r
				}

				if dur, _, err := parseLapStr(data["avg"]); err == nil {
					lt.ClassFT.AvgLapDuration = dur
				}
			} else {
				slog.Warn("failed to parse Class FT meta", "raw", text)
			}
		case strings.HasPrefix(text, "Class Best Lap:"):
			if data := namedCapture(classBestLapRegex, text); data != nil {
				lt.ClassBestLap = &models.ClassBestLap{DriverName: ptr(data["name"])}
				dur, _, _ := parseLapStr(data["time"])
				lt.ClassBestLap.Time = dur
			} else {
				slog.Warn("failed to parse Class Best Lap meta", "raw", text)
			}
		}
	})
}
