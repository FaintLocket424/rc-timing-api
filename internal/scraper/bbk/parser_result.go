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
	practiceRegex     = regexp.MustCompile(`Practice (?P<practice>\d+) \((?P<class>.*?)\)(?: R(?P<round>\d+))?`)
	heatRegex         = regexp.MustCompile(`Heat (?P<heat>\d+) \((?P<class>.*?)\)(?: R(?P<round>\d+))?`)
	finalRegex        = regexp.MustCompile(`(?P<final>[A-Z])\. Final \((?P<class>.*?)\)(?: L(?P<leg>\d+))?`)
	bestLapRegex      = regexp.MustCompile(`Best Lap: (?P<name>.*?) (?P<time>[\d\.]+) L:(?P<lap>\d+)`)
	classFTRegex      = regexp.MustCompile(`Class FT: (?P<name>.*?) (?P<res>[\d\/'\.]+) Av Lap (?P<avg>[\d\.]+)(?: R:(?P<round>\d+))?`)
	classBestLapRegex = regexp.MustCompile(`Class Best Lap: (?P<name>.*?) (?P<time>\d+\.\d+)`)
)

func parseRaceResultHTML(body io.Reader) (*models.RaceResultScrape, error) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML body: %w", err)
	}

	tables := doc.Find("table.NBT")
	if tables.Length() < 3 {
		return nil, fmt.Errorf("critical structure change: expected 3 tables, found %d", tables.Length())
	}

	scrape := &models.RaceResultScrape{}

	parseHeader(tables.Eq(0), scrape)
	parseDrivers(tables.Eq(1), scrape)
	parseMeta(tables.Eq(2), scrape)

	return scrape, nil
}

func parseHeader(s *goquery.Selection, lt *models.RaceResultScrape) {
	tds := s.Find("td.livehtml-text")
	if tds.Length() < 2 {
		slog.Warn("header table malformed, skipping header details")
		return
	}

	text := tds.First().Text()

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

	timerText := tds.Last().Text()
	if before, _, ok := strings.Cut(timerText, "("); ok {
		var m, s int
		if _, err := fmt.Sscanf(before, "%d'%ds", &m, &s); err == nil {
			lt.ElapsedTime = ptr(time.Duration(m)*time.Minute + time.Duration(s)*time.Second)
		} else {
			slog.Warn("could not parse timer", "raw", timerText)
		}
	}
}

func parseDrivers(drivers *goquery.Selection, lt *models.RaceResultScrape) {
	rows := drivers.Find("tbody tr")
	if rows.Length() < 2 {
		slog.Error("drivers table missing data rows")
		return
	}

	// Map headers once
	idx := make(map[string]int)
	rows.First().Find("td").Each(func(i int, s *goquery.Selection) {
		text := s.Text()
		switch text {
		case "C":
			idx["car"] = i
		case "Driver":
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
	rows.Slice(1, rows.Length()).Each(func(i int, s *goquery.Selection) {
		cols := s.Find("td")
		dr := models.DriverRaceResult{}

		dr.Name = ptr(cols.Eq(idx["name"]).Text())
		if v, err := strconv.Atoi(cols.Eq(idx["car"]).Text()); err == nil {
			dr.CarNumber = ptr(v)
		}

		// Result Parsing
		resText := cols.Eq(idx["res"]).Text()

		if laps, dur, err := parseRaceResultStr(resText); err == nil {
			dr.Laps = laps
			if *dr.Laps < 0 && len(lt.Drivers) > 0 && lt.Drivers[0].Laps != nil {
				*dr.Laps += *lt.Drivers[0].Laps
			}
			dr.Time = dur
		} else if gap, err := parseGapStr(resText); err == nil {
			if len(lt.Drivers) > 0 && lt.Drivers[0].Laps != nil && lt.Drivers[0].Time != nil {
				dr.Laps = ptr(*lt.Drivers[0].Laps)
				dr.Time = ptr(*lt.Drivers[0].Time + *gap)
			}
			// } else if resText == "" {
			// 	// empty string may mean DNS
			// } else if strings.HasPrefix(resText, "W-") {
			// 	// warm up laps
			// } else if resText == "DNS" {
			// 	// Did not start
		} else {
			slog.Warn("unparseable result", "row", i+1, "result", resText)
		}

		if dur, ln, err := parseLapStr(cols.Eq(idx["best"]).Text()); err == nil {
			dr.BestLapDuration = dur
			dr.BestLapNumber = ln
		}

		if dur, _, err := parseLapStr(cols.Eq(idx["last"]).Text()); err == nil {
			dr.LastLapDuration = dur
		}

		lt.Drivers = append(lt.Drivers, dr)
	})
}

func parseMeta(meta *goquery.Selection, lt *models.RaceResultScrape) {
	meta.Find("tr td").Each(func(_ int, s *goquery.Selection) {
		font := s.Find("font").First().Text()
		text := s.Text()

		switch font {
		case "Best Lap:":
			if data := NamedCapture(bestLapRegex, text); data != nil {
				lt.BestLap = &models.RaceBestLap{DriverName: ptr(data["name"])}
				dur, ln, _ := parseLapStr(fmt.Sprintf("%s[%s]", data["time"], data["lap"]))
				lt.BestLap.Time, lt.BestLap.LapNumber = dur, ln
			} else {
				slog.Warn("failed to parse Best Lap meta", "raw", text)
			}
		case "Class FT:":
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
		case "Class Best Lap:":
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
