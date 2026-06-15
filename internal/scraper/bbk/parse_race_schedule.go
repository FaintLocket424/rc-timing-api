package bbk

import (
	"fmt"
	"io"
	"log/slog"
	"strings"
	"time"

	"github.com/FaintLocket424/opengrid-bridge/internal/models"
	"github.com/PuerkitoBio/goquery"
)

func parseRaceScheduleHTML(body io.Reader) (*models.RaceScheduleScrape, error) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML body: %w", err)
	}

	scrape := &models.RaceScheduleScrape{
		Practice:   []models.ScheduledRace{},
		Qualifying: []models.ScheduledRace{},
		Finals:     []models.ScheduledRace{},
	}
	scrape.Title, scrape.Timestamp = parseTitleAndTimestamp(doc)

	// Dynamically investigate tables to locate the schedule grid.
	var scheduleTable *goquery.Selection
	doc.Find("table").Each(func(_ int, t *goquery.Selection) {
		if scheduleTable != nil {
			return
		}
		t.Find("td").Each(func(_ int, td *goquery.Selection) {
			if strings.TrimSpace(td.Text()) == "Start Time" {
				scheduleTable = t
			}
		})
	})

	if scheduleTable != nil {
		parseScheduledRacesTable(scheduleTable, scrape)
	}

	return scrape, nil
}

func parseScheduledRacesTable(table *goquery.Selection, lt *models.RaceScheduleScrape) {
	rows := table.Find("tr")
	if rows.Length() < 2 {
		slog.Error("schedule table missing data rows")
		return
	}

	headerRow, headerIdx, found := findHeaderRow(rows, "Start Time")
	if !found {
		slog.Error("schedule table missing column headers")
		return
	}

	idx := mapColumnHeaders(headerRow, map[string][]string{
		"race": {"Race"},
		"time": {"Start Time"},
	})

	rows.Slice(headerIdx+1, rows.Length()).Each(func(_ int, s *goquery.Selection) {
		cols := s.Children()

		var p, h, f, r *int
		var class *string
		var matched bool

		if raceIdx, ok := idx["race"]; ok && raceIdx < cols.Length() {
			p, h, f, r, class, matched = parseRaceTitleText(cols.Eq(raceIdx).Text())
		}

		if !matched {
			return
		}

		var startTime *time.Time
		if startTimeIdx, ok := idx["time"]; ok && startTimeIdx < cols.Length() {
			if t, err := parseTimeToday(cols.Eq(startTimeIdx).Text()); err == nil {
				startTime = t
			}
		}

		roundVal := 1
		if r != nil {
			roundVal = *r
		}

		switch {
		case p != nil:
			sr := models.ScheduledRace{
				HeatRound: &models.HeatRound{Heat: *p, Round: roundVal},
				ClassName: class,
				StartTime: startTime,
			}
			lt.Practice = append(lt.Practice, sr)
		case h != nil:
			sr := models.ScheduledRace{
				HeatRound: &models.HeatRound{Heat: *h, Round: roundVal},
				ClassName: class,
				StartTime: startTime,
			}
			lt.Qualifying = append(lt.Qualifying, sr)
		case f != nil:
			sr := models.ScheduledRace{
				HeatRound: &models.HeatRound{Heat: *f, Round: roundVal},
				ClassName: class,
				StartTime: startTime,
			}
			lt.Finals = append(lt.Finals, sr)
		}
	})
}
