package bbk

import (
	"fmt"
	"io"
	"regexp"
	"strconv"
	"time"

	"github.com/FaintLocket424/rc-timing-api/internal/models"
	"github.com/FaintLocket424/rc-timing-api/internal/scraper/bbk/utils"
	"github.com/PuerkitoBio/goquery"
)

var resultsLinkRegex = regexp.MustCompile(`^(?P<type>[phf])(?P<heat>\d+)r(?P<round>\d+)res\.htm$`)

func parseRaceResultsIndex(body io.Reader) (*models.RaceResultsIndexScrape, error) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML body: %w", err)
	}

	tables := doc.Find("table.NBT")
	if l := tables.Length(); l < 2 {
		return nil, fmt.Errorf("expected at least 2 tables, found %d", l)
	}

	scrape := &models.RaceResultsIndexScrape{
		Practice: make(map[models.HeatRound]struct{}),
		Quali:    make(map[models.HeatRound]struct{}),
		Finals:   make(map[models.HeatRound]struct{}),
	}

	header := tables.First().Find("td")
	scrape.Title = header.First().Text()

	if t, err := time.Parse("15:04", header.Last().Text()); err != nil {
		now := time.Now()

		scrape.Timestamp = time.Date(now.Year(), now.Month(), now.Day(), t.Hour(), t.Minute(), 0, 0, now.Location())
	}

	links := tables.Find("a")

	links.Each(func(i int, s *goquery.Selection) {
		if attr, exists := s.Attr("href"); exists {
			if data := utils.NamedCapture(resultsLinkRegex, attr); data != nil {
				heat, err := strconv.Atoi(data["heat"])
				if err != nil {
					return
				}

				round, err := strconv.Atoi(data["round"])
				if err != nil {
					return
				}

				switch data["type"] {
				case "p":
					scrape.Practice[models.HeatRound{Heat: heat, Round: round}] = struct{}{}
				case "h":
					scrape.Quali[models.HeatRound{Heat: heat, Round: round}] = struct{}{}
				case "f":
					scrape.Finals[models.HeatRound{Heat: heat, Round: round}] = struct{}{}
				}
			}
		}
	})

	return scrape, nil
}
