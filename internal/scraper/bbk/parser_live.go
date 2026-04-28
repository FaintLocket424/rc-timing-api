package bbk

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/FaintLocket424/rc-timing-api/internal/models"
	"github.com/PuerkitoBio/goquery"
)

var (
	heatRegex          = regexp.MustCompile(`Heat (?P<heat>\d+) \((?P<class>.*?)\)(?: R(?P<round>\d+))?`)
	finalRegex         = regexp.MustCompile(`(?P<round>[A-Z])\. Final \((?P<class>.*?)\)(?: L(?P<leg>\d+))?`)
	driversResultRegex = regexp.MustCompile(`(?P<laps>-?\d+)/(?:(?P<mins>\d+)')?(?P<secs>\d+\.\d+)`)
	lapFormatRegex     = regexp.MustCompile(`(?P<time>\d+\.\d+)(?:\[(?P<lap>\d+)\])?`)
	bestLapRegex       = regexp.MustCompile(`Best Lap: (?P<name>.*?) (?P<time>[\d\.]+) L:(?P<lap>\d+)`)
	classFTRegex       = regexp.MustCompile(`Class FT: (?P<name>.*?) (?P<laps>\d+)/(?P<time>[\d\.\']+) Av Lap (?P<avg>[\d\.]+)(?: R:(?P<round>\d+))?`)
	classBestLapRegex  = regexp.MustCompile(`Class Best Lap: (?P<name>.*?) (?P<time>\d+\.\d+)`)
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

func (s *BBKScraper) GetLiveTiming() (*models.LiveTimingScrape, error) {
	res, err := s.Client.Get(s.Target + "/liveraceres.htm")
	if err != nil {
		return nil, fmt.Errorf("network error: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned status: %d", res.StatusCode)
	}

	return parseRaceResult(res.Body)
}

func parseRaceResult(body io.Reader) (*models.LiveTimingScrape, error) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML body: %w", err)
	}

	tables := doc.Find("table.NBT")
	if tables.Length() < 3 {
		return nil, fmt.Errorf("critical structure change: expected 3 tables, found %d", tables.Length())
	}

	lt := &models.LiveTimingScrape{}
	parseHeader(tables.Eq(0), lt)
	parseDrivers(tables.Eq(1), lt)
	parseMeta(tables.Eq(2), lt)

	return lt, nil
}

func parseHeader(s *goquery.Selection, lt *models.LiveTimingScrape) {
	tds := s.Find("td.livehtml-text")
	if tds.Length() < 2 {
		slog.Warn("header table malformed, skipping header details")
		return
	}

	text := tds.First().Text()

	if data := NamedCapture(heatRegex, text); data != nil {
		lt.HeatNumber, _ = strconv.Atoi(data["heat"])
		lt.ClassName = data["class"]
		lt.Round, _ = strconv.Atoi(data["round"])
	} else if data := NamedCapture(finalRegex, text); data != nil {
		lt.ClassName = data["class"]

		if data["leg"] != "" {
			lt.Round, _ = strconv.Atoi(data["leg"])
		} else {
			lt.Round = 1
		}
	}

	// Timer
	timerText := tds.Last().Text()
	if before, _, ok := strings.Cut(timerText, "("); ok {
		var m, s int
		if _, err := fmt.Sscanf(before, "%d'%ds", &m, &s); err == nil {
			lt.ElapsedTime = time.Duration(m)*time.Minute + time.Duration(s)*time.Second
		} else {
			slog.Warn("could not parse timer", "raw", timerText)
		}
	}
}

func parseMillis(s string) (time.Duration, error) {
	parts := strings.Split(s, ".")

	// Parse whole seconds
	secs, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, err
	}

	// Handle fractional part (milliseconds)
	millis := 0
	if len(parts) > 1 {
		frac := parts[1]
		// Pad or truncate to 3 digits (e.g., "23" -> 230ms, "2345" -> 234ms)
		for len(frac) < 3 {
			frac += "0"
		}
		millis, err = strconv.Atoi(frac[:3])
		if err != nil {
			return 0, err
		}
	}

	return time.Duration(secs)*time.Second + time.Duration(millis)*time.Millisecond, nil
}

func parseDrivers(drivers *goquery.Selection, lt *models.LiveTimingScrape) {
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

		dr.Name = cols.Eq(idx["name"]).Text()
		dr.CarNumber, _ = strconv.Atoi(cols.Eq(idx["car"]).Text())

		// Result Parsing
		resText := cols.Eq(idx["res"]).Text()

		if data := NamedCapture(driversResultRegex, resText); data != nil {
			dr.Laps, _ = strconv.Atoi(data["laps"])
			if dr.Laps < 0 && len(lt.Drivers) > 0 {
				dr.Laps += lt.Drivers[0].Laps
			}

			var totalDuration time.Duration

			// 1. Add Minutes
			if data["mins"] != "" {
				mins, _ := strconv.Atoi(data["mins"])
				totalDuration += time.Duration(mins) * time.Minute
			}

			// 2. Add Seconds + Millis
			secondsPart, _ := parseMillis(data["secs"])
			totalDuration += secondsPart

			dr.Time = totalDuration
		} else if strings.HasPrefix(resText, "+") {
			gapDuration, err := parseMillis(resText[1:]) // Pass the string after '+'
			if err == nil && len(lt.Drivers) > 0 {
				leader := lt.Drivers[0]
				dr.Laps = leader.Laps
				dr.Time = leader.Time + gapDuration
			}
		} else if resText == "" {
			// empty string may mean DNS
		} else if strings.HasPrefix(resText, "W-") {
			// warm up laps
		} else {
			slog.Warn("unparseable result", "row", i+1, "result", resText)
		}

		if bestLapData := NamedCapture(lapFormatRegex, cols.Eq(idx["best"]).Text()); bestLapData != nil {
			dr.BestLapDuration, _ = time.ParseDuration(bestLapData["time"] + "s")
			dr.BestLapNumber, _ = strconv.Atoi(bestLapData["lap"])
		}

		if lastLapData := NamedCapture(lapFormatRegex, cols.Eq(idx["last"]).Text()); lastLapData != nil {
			dr.LastLapDuration, _ = time.ParseDuration(lastLapData["time"] + "s")
		}

		lt.Drivers = append(lt.Drivers, dr)
	})
}

func parseMeta(meta *goquery.Selection, lt *models.LiveTimingScrape) {
	meta.Find("tr td").Each(func(i int, s *goquery.Selection) {
		font := s.Find("font").First().Text()
		text := s.Text()

		switch font {
		case "Best Lap:":
			if data := NamedCapture(bestLapRegex, text); data != nil {
				lt.BestLap.DriverName = data["name"]
				lt.BestLap.Time, _ = time.ParseDuration(data["time"] + "s")
				lt.BestLap.LapNumber, _ = strconv.Atoi(data["lap"])
			} else {
				slog.Warn("failed to parse Best Lap meta", "raw", text)
			}
		case "Class FT:":
			if data := NamedCapture(classFTRegex, text); data != nil {
				lt.ClassFT.DriverName = data["name"]
				lt.ClassFT.Laps, _ = strconv.Atoi(data["laps"])
				lt.ClassFT.Time, _ = time.ParseDuration(strings.Replace(data["time"], "'", "m", 1) + "s")
				lt.ClassFT.AvgLapDuration, _ = time.ParseDuration(data["avg"] + "s")
				lt.ClassFT.Round, _ = strconv.Atoi(data["round"])
			} else {
				slog.Warn("failed to parse Class FT meta", "raw", text)
			}
		case "Class Best Lap:":
			if data := NamedCapture(classBestLapRegex, text); data != nil {
				lt.ClassBestLap.DriverName = data["name"]
				lt.ClassBestLap.Time, _ = time.ParseDuration(data["time"] + "s")
			} else {
				slog.Warn("failed to parse Class Best Lap meta", "raw", text)
			}
		}

	})
}
