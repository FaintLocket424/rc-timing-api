package bbk

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/FaintLocket424/rc-timing-api/internal/models"
	"github.com/PuerkitoBio/goquery"
	"github.com/fatih/color"
)

var red = color.New(color.FgRed)

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

func (s BBKScraper) GetLiveTiming() (*models.LiveTimingScrape, error) {
	res, err := s.Client.Get(s.Target + "/liveraceres.htm")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code error: %d, %s", res.StatusCode, res.Status)
	}

	return parseRaceResult(res.Body)
}

func parseRaceResult(body io.Reader) (*models.LiveTimingScrape, error) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse document: %w", err)
	}

	tables := doc.Find("table.NBT")

	expectedTableCount := 3
	if tables.Length() != expectedTableCount {
		return nil, fmt.Errorf("expected %d tables, found %s", expectedTableCount, red.Sprint(tables.Length()))
	}

	lt := &models.LiveTimingScrape{}

	if err := parseHeader(tables.Eq(0), lt); err != nil {
		return nil, fmt.Errorf("failed to parse header: %w", err)
	}

	if err := parseDrivers(tables.Eq(1), lt); err != nil {
		return nil, fmt.Errorf("failed to parse drivers: %w", err)
	}

	if err := parseMeta(tables.Eq(2), lt); err != nil {
		return nil, fmt.Errorf("failed to parse meta: %w", err)
	}

	return lt, nil
}

func parseHeader(header *goquery.Selection, lt *models.LiveTimingScrape) error {
	tds := header.Find("td.livehtml-text")

	if tds.Length() != 2 {
		return fmt.Errorf("expected 2 td elements, found %s", red.Sprint(tds.Length()))
	}

	heat, class, round, err := parseHeaderTitle(tds.First().Text())
	if err != nil {
		return fmt.Errorf("failed to parse header title: %w", err)
	}
	lt.HeatNumber, lt.ClassName, lt.Round = heat, class, round

	elapsed, err := parseHeaderTimer(tds.Last().Text())
	if err != nil {
		return fmt.Errorf("failed to parse header timer: %w", err)
	}
	lt.ElapsedTime = elapsed

	return nil
}

var headerTitleRegex = regexp.MustCompile(`Heat (?P<heat>\d+) \((?P<class>.*?)\) R(?P<round>\d+)`)

func parseHeaderTitle(input string) (int, string, int, error) {
	data := NamedCapture(headerTitleRegex, input)
	if data == nil {
		return 0, "", 0, fmt.Errorf("failed to parse heat header: %q", input)
	}

	heatNum, err := strconv.Atoi(data["heat"])
	if err != nil {
		return 0, "", 0, fmt.Errorf("invalid heat number %q: %w", data["heat"], err)
	}

	round, err := strconv.Atoi(data["round"])
	if err != nil {
		return 0, "", 0, fmt.Errorf("invalid round number %q: %w", data["round"], err)
	}

	return heatNum, data["class"], round, nil
}

func parseHeaderTimer(input string) (time.Duration, error) {
	bracketIdx := strings.Index(input, "(")
	if bracketIdx == -1 {
		return 0, fmt.Errorf("unexpected timer format: %q", input)
	}

	elapsedStr := strings.TrimSpace(input[:bracketIdx])

	var minutes, seconds int
	_, err := fmt.Sscanf(elapsedStr, "%d'%ds", &minutes, &seconds)
	if err != nil {
		return 0, fmt.Errorf("failed to parse duration %q: %w", elapsedStr, err)
	}

	duration := time.Duration(minutes)*time.Minute + time.Duration(seconds)*time.Second
	return duration, nil
}

var driversResultRegex = regexp.MustCompile(`(?P<laps>-?\d+)/(?P<mins>\d+)'(?P<secs>\d+\.\d+)`)

func parseDrivers(drivers *goquery.Selection, lt *models.LiveTimingScrape) error {
	tableRows := drivers.Find("tbody tr")
	if tableRows.Length() < 2 {
		return fmt.Errorf("insufficient rows in drivers table")
	}

	var carNumIdx, driverNameIdx, resultIdx, bestLapIdx = -1, -1, -1, -1
	tableRows.First().Find("td").Each(func(i int, s *goquery.Selection) {
		switch s.Text() {
		case "C":
			carNumIdx = i
		case "Driver":
			driverNameIdx = i
		case "Result":
			resultIdx = i
		case "B-Lp":
			bestLapIdx = i
		}
	})

	var errReturn error
	tableRows.Slice(1, tableRows.Length()).EachWithBreak(func(i int, s *goquery.Selection) bool {
		cols := s.Find("td")
		dr := models.DriverRaceResult{}
		defer func() {
			lt.Drivers = append(lt.Drivers, dr)
		}()

		val := strings.TrimSpace(cols.Eq(carNumIdx).Text())
		if num, err := strconv.Atoi(val); err == nil {
			dr.CarNumber = num
		} else {
			errReturn = fmt.Errorf("row %d: invalid car number %q: %w", i+1, val, err)
			return false
		}

		dr.Name = strings.TrimSpace(cols.Eq(driverNameIdx).Text())

		resText := cols.Eq(resultIdx).Text()
		data := NamedCapture(driversResultRegex, resText)
		if data == nil {
			// errReturn = fmt.Errorf("row %d: failed to parse result format %q", i+1, resText)
			// return false
			return true
		}

		laps, err := strconv.Atoi(data["laps"])
		if err != nil {
			errReturn = fmt.Errorf("row %d: invalid laps %q: %w", i+1, data["laps"], err)
			return false
		}

		if laps < 0 {
			laps += lt.Drivers[0].Laps
		}
		dr.Laps = laps

		mins, err := strconv.Atoi(data["mins"])
		if err != nil {
			errReturn = fmt.Errorf("row %d: invalid minutes %q: %w", i+1, data["mins"], err)
			return false
		}

		secs, err := strconv.ParseFloat(data["secs"], 64)
		if err != nil {
			errReturn = fmt.Errorf("row %d: invalid seconds %q: %w", i+1, data["secs"], err)
			return false
		}
		dr.Time = time.Duration(mins)*time.Minute + time.Duration(secs*float64(time.Second))

		blStr := strings.TrimSpace(cols.Eq(bestLapIdx).Text())
		if dur, err := time.ParseDuration(blStr + "s"); err == nil {
			dr.BestLapDuration = dur
		} else {
			errReturn = fmt.Errorf("row %d: invalid best lap %q: %w", i+1, blStr, err)
			return false
		}

		return true
	})

	return errReturn
}

var bestLapRegex = regexp.MustCompile(`Best Lap: (?P<name>.*?) (?P<time>[\d\.]+) L:(?P<lap>\d+)`)

func parseMeta(meta *goquery.Selection, lt *models.LiveTimingScrape) error {
	lines := meta.Find("tr td.fastest-lap")

	if lines.Length() != 3 {
		return fmt.Errorf("expected 3 metadata cells, found %d", lines.Length())
	}

	if err := parseMetaBestLap(lines.Eq(0).Text(), lt); err != nil {
		return fmt.Errorf("failed to parse best lap: %w", err)
	}

	if err := parseMetaClassFT(lines.Eq(1).Text(), lt); err != nil {
		return fmt.Errorf("failed to parse class FT: %w", err)
	}

	if err := parseMetaClassBestLap(lines.Eq(2).Text(), lt); err != nil {
		return fmt.Errorf("failed to parse class best lap: %w", err)
	}

	return nil
}

func parseMetaBestLap(input string, lt *models.LiveTimingScrape) error {
	data := NamedCapture(bestLapRegex, input)
	if data == nil {
		return fmt.Errorf("unexpected best lap format: %q", input)
	}

	lt.BestLap.DriverName = data["name"]

	lapDuration, err := time.ParseDuration(data["time"] + "s")
	if err != nil {
		return fmt.Errorf("invalid duration in %q: %w", input, err)
	}
	lt.BestLap.Time = lapDuration

	lapNumber, err := strconv.Atoi(data["lap"])
	if err != nil {
		return fmt.Errorf("invalid lap number in %q: %w", input, err)
	}
	lt.BestLap.LapNumber = lapNumber

	return nil
}

var classFTRegex = regexp.MustCompile(`Class FT: (?P<name>.*?) (?P<laps>\d+)/(?P<time>[\d\.\']+) Av Lap (?P<avg>[\d\.]+) R:(?P<round>\d+)`)

func parseMetaClassFT(input string, lt *models.LiveTimingScrape) error {
	data := NamedCapture(classFTRegex, input)
	if data == nil {
		return fmt.Errorf("unexpected Class FT format: %q", input)
	}

	var err error

	lt.ClassFT.DriverName = data["name"]

	lt.ClassFT.Laps, err = strconv.Atoi(data["laps"])
	if err != nil {
		return fmt.Errorf("invalid laps in %q: %w", input, err)
	}

	// Parse duration: format "M'SS.s" -> "MmSS.ss"
	formattedTime := strings.Replace(data["time"], "'", "m", 1) + "s"
	lt.ClassFT.Time, err = time.ParseDuration(formattedTime)
	if err != nil {
		return fmt.Errorf("invalid time in %q: %w", input, err)
	}

	lt.ClassFT.AvgLapDuration, err = time.ParseDuration(data["avg"] + "s")
	if err != nil {
		return fmt.Errorf("invalid avg lap in %q: %w", input, err)
	}

	lt.ClassFT.Round, err = strconv.Atoi(data["round"])
	if err != nil {
		return fmt.Errorf("invalid round in %q: %w", input, err)
	}

	return nil
}

var classBestLapRegex = regexp.MustCompile(`Class Best Lap: (?P<name>.*?) (?P<time>\d+\.\d+)`)

func parseMetaClassBestLap(input string, lt *models.LiveTimingScrape) error {
	data := NamedCapture(classBestLapRegex, input)
	if data == nil {
		return fmt.Errorf("unexpected Class Best Lap format: %q", input)
	}

	lt.ClassBestLap.DriverName = data["name"]

	d, err := time.ParseDuration(data["time"] + "s")
	if err != nil {
		return fmt.Errorf("invalid Class Best Lap duration in %q: %w", input, err)
	}
	lt.ClassBestLap.Time = d

	return nil
}
