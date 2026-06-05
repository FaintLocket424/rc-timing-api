package bbk

import (
	"fmt"
	"io"
	"log/slog"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/FaintLocket424/opengrid-bridge/internal/models"
	"github.com/PuerkitoBio/goquery"
)

var (
	practiceRegex           = regexp.MustCompile(`Practice\s+(?P<practice>\d+)\s+\((?P<class>.*?)\)(?:\s+R(?P<round>\d+))?`)
	heatRegex               = regexp.MustCompile(`Heat\s+(?P<heat>\d+)\s+\((?P<class>.*?)\)(?:\s+R(?P<round>\d+))?`)
	finalRegex              = regexp.MustCompile(`(?P<final>[A-Z])\.\s+Final\s+\((?P<class>.*?)\)(?:\s+L(?P<leg>\d+))?`)
	bestLapRegex            = regexp.MustCompile(`Best Lap:\s*(?P<name>.*?)\s+(?P<time>[\d\.]+)\s+L\s*:\s*(?P<lap>\d+)`)
	classFTRegex            = regexp.MustCompile(`Class FT:\s*(?P<name>.*?)\s+(?P<res>[\d\/'\.]+)\s+Av Lap\s+(?P<avg>[\d\.]+)(?:\s+R\s*:\s*(?P<round>\d+))?`)
	classBestLapRegex       = regexp.MustCompile(`Class Best Lap:\s*(?P<name>.*?)\s+(?P<time>\d+\.\d+)`)
	finishedTimeHeaderRegex = regexp.MustCompile(`Finished\s+(?P<finishTime>\d{2}:\d{2})`)
	elapsedTimeHeaderRegex  = regexp.MustCompile(`(?P<elapsed>\d+'\d+s)\s*\((?P<remaining>\d+s)\)`)
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

	return scrape, nil
}

func parseHeader(s *goquery.Selection, lt *models.RaceResultScrape) {
	s.Find("td").Each(func(_ int, td *goquery.Selection) {
		text := strings.TrimSpace(td.Text())

		if data := NamedCapture(practiceRegex, text); data != nil {
			if v, err := strconv.Atoi(data["practice"]); err == nil {
				lt.PracticeNumber = ptr(v)
			}
			lt.ClassName = ptr(data["class"])
			if v, err := strconv.Atoi(data["round"]); err == nil {
				lt.Round = ptr(v)
			}
		} else if data := NamedCapture(heatRegex, text); data != nil {
			if v, err := strconv.Atoi(data["heat"]); err == nil {
				lt.HeatNumber = ptr(v)
			}
			lt.ClassName = ptr(data["class"])
			if v, err := strconv.Atoi(data["round"]); err == nil {
				lt.Round = ptr(v)
			}
		} else if data := NamedCapture(finalRegex, text); data != nil {
			final := data["final"]

			if len(final) == 1 {
				r := strings.ToUpper(final)[0]
				v := int(r - 'A' + 1)
				lt.FinalNumber = ptr(v)
			}

			lt.ClassName = ptr(data["class"])

			if data["leg"] != "" {
				if v, err := strconv.Atoi(data["leg"]); err == nil {
					lt.Round = ptr(v)
				}
			} else {
				lt.Round = ptr(1)
			}
		}

		if finishMatch := NamedCapture(finishedTimeHeaderRegex, text); finishMatch != nil {
			t, err := time.Parse("15:04", finishMatch["finishTime"])
			if err == nil {
				now := time.Now()
				lt.FinishTime = ptr(time.Date(now.Year(), now.Month(), now.Day(), t.Hour(), t.Minute(), 0, 0, now.Location()))
			}
			lt.RaceStatus = ptr("Finished")
		}

		if durMatch := NamedCapture(elapsedTimeHeaderRegex, text); durMatch != nil {
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

	// Dynamically scan for the column headers row (which contains "Driver")
	var headerRow *goquery.Selection
	var headerIdx int
	rows.Each(func(i int, r *goquery.Selection) {
		if headerRow != nil {
			return
		}
		r.Children().Each(func(_ int, c *goquery.Selection) {
			if strings.TrimSpace(c.Text()) == driverColName {
				headerRow = r
				headerIdx = i
			}
		})
	})

	if headerRow == nil {
		slog.Error("drivers table missing column headers")
		return
	}

	// Map headers
	idx := make(map[string]int)
	headerRow.Children().Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		switch text {
		case "C":
			idx["car"] = i
		case driverColName:
			idx["name"] = i
		case "Result":
			idx["res"] = i
		case "B-Lap", "B-Lp":
			idx["best"] = i
		case "L-Lap", "L-Lp":
			idx["last"] = i
		}
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
			if v, err := strconv.Atoi(strings.TrimSpace(cols.Eq(carIdx).Text())); err == nil {
				dr.CarNumber = ptr(v)
			}
		}

		// Result
		if resIdx, ok := idx["res"]; ok && resIdx < cols.Length() {
			resText := strings.TrimSpace(cols.Eq(resIdx).Text())

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
				i, err := strconv.Atoi(numStr)
				if err == nil {
					dr.WarmupLaps = ptr(i)
				}
			} else if resText != "" {
				slog.Warn("unparseable result", "row", i+1, "result", resText)
			}
		}

		// Best Lap duration
		if bestIdx, ok := idx["best"]; ok && bestIdx < cols.Length() {
			bestText := strings.TrimSpace(cols.Eq(bestIdx).Text())
			if dur, ln, err := parseLapStr(bestText); err == nil {
				dr.BestLapDuration = dur
				dr.BestLapNumber = ln
			}
		}

		// Last Lap duration
		if lastIdx, ok := idx["last"]; ok && lastIdx < cols.Length() {
			lastText := strings.TrimSpace(cols.Eq(lastIdx).Text())
			if dur, _, err := parseLapStr(lastText); err == nil {
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
			if data := NamedCapture(bestLapRegex, text); data != nil {
				lt.BestLap = &models.RaceBestLap{DriverName: ptr(data["name"])}
				dur, ln, _ := parseLapStr(fmt.Sprintf("%s[%s]", data["time"], data["lap"]))
				lt.BestLap.Time, lt.BestLap.LapNumber = dur, ln
			} else {
				slog.Warn("failed to parse Best Lap meta", "raw", text)
			}
		case strings.HasPrefix(text, "Class FT:"):
			if data := NamedCapture(classFTRegex, text); data != nil {
				lt.ClassFT = &models.ClassFT{DriverName: ptr(data["name"])}
				laps, dur, _ := parseRaceResultStr(data["res"])
				lt.ClassFT.Laps, lt.ClassFT.Time = laps, dur
				if r, err := strconv.Atoi(data["round"]); err == nil {
					lt.ClassFT.Round = ptr(r)
				}
			} else {
				slog.Warn("failed to parse Class FT meta", "raw", text)
			}
		case strings.HasPrefix(text, "Class Best Lap:"):
			if data := NamedCapture(classBestLapRegex, text); data != nil {
				lt.ClassBestLap = &models.ClassBestLap{DriverName: ptr(data["name"])}
				dur, _, _ := parseLapStr(data["time"])
				lt.ClassBestLap.Time = dur
			} else {
				slog.Warn("failed to parse Class Best Lap meta", "raw", text)
			}
		}
	})
}
