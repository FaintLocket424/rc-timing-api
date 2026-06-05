package bbk

import (
	"cmp"
	"fmt"
	"io"
	"regexp"
	"slices"
	"strconv"
	"time"

	"github.com/FaintLocket424/opengrid-bridge/internal/models"
	"github.com/PuerkitoBio/goquery"
)

var (
	resultsLinkRegex       = regexp.MustCompile(`^(?P<type>[phf])(?P<heat>\d+)r(?P<round>\d+)res\.htm$`)
	roundOverallsLinkRegex = regexp.MustCompile(`^(?P<type>[phf])eor(?P<round>\d+)\.htm$`)
)

func parseRaceResultsIndexHTML(body io.Reader) (*models.RaceResultsIndexScrape, error) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML body: %w", err)
	}

	tables := doc.Find("table.NBT")
	if l := tables.Length(); l < 2 {
		return nil, fmt.Errorf("expected at least 2 tables, found %d", l)
	}

	scrape := &models.RaceResultsIndexScrape{
		Practice: models.ResultStatus{
			Results:       []models.HeatRound{},
			RoundOveralls: []int{},
		},
		Qualifying: models.ResultStatus{
			Results:       []models.HeatRound{},
			RoundOveralls: []int{},
		},
		Finals: models.ResultStatus{
			Results:       []models.HeatRound{},
			RoundOveralls: []int{},
		},
	}

	header := tables.First().Find("td")

	if header.Length() > 1 {
		if text := header.First().Text(); text != "" {
			scrape.Title = ptr(text)
		}
	}

	if t, err := time.Parse("15:04", header.Last().Text()); err == nil {
		now := time.Now()

		scrape.Timestamp = time.Date(now.Year(), now.Month(), now.Day(), t.Hour(), t.Minute(), 0, 0, now.Location())
	}

	links := tables.Find("a")

	links.Each(func(_ int, s *goquery.Selection) {
		if attr, exists := s.Attr("href"); exists {
			if data := NamedCapture(resultsLinkRegex, attr); data != nil {
				heat, err := strconv.Atoi(data["heat"])
				if err != nil {
					return
				}

				round, err := strconv.Atoi(data["round"])
				if err != nil {
					return
				}

				hr := models.HeatRound{Heat: heat, Round: round}

				switch data["type"] {
				case "p":
					scrape.Practice.Results = append(scrape.Practice.Results, hr)
				case "h":
					scrape.Qualifying.Results = append(scrape.Qualifying.Results, hr)
				case "f":
					scrape.Finals.Results = append(scrape.Finals.Results, hr)
				}
			} else if data := NamedCapture(roundOverallsLinkRegex, attr); data != nil {
				round, err := strconv.Atoi(data["round"])
				if err != nil {
					return
				}

				switch data["type"] {
				case "p":
					if round > 0 {
						scrape.Practice.RoundOveralls = append(scrape.Practice.RoundOveralls, round)
					} else {
						scrape.Practice.Overall = true
					}
				case "h":
					if round > 0 {
						scrape.Qualifying.RoundOveralls = append(scrape.Qualifying.RoundOveralls, round)
					} else {
						scrape.Qualifying.Overall = true
					}
				case "f":
					if round > 0 {
						scrape.Finals.RoundOveralls = append(scrape.Finals.RoundOveralls, round)
					} else {
						scrape.Finals.Overall = true
					}
				}
			}
		}
	})

	// Sorts by Round ascending, then Heat ascending
	sortFunc := func(a, b models.HeatRound) int {
		if a.Round != b.Round {
			return cmp.Compare(a.Round, b.Round)
		}
		return cmp.Compare(a.Heat, b.Heat)
	}

	slices.SortFunc(scrape.Practice.Results, sortFunc)
	slices.SortFunc(scrape.Qualifying.Results, sortFunc)
	slices.SortFunc(scrape.Finals.Results, sortFunc)

	return scrape, nil
}
