package bbk

import (
	"fmt"
	"io"
	"regexp"
	"strconv"

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

	scrape := &models.RaceResultsIndexScrape{}

	// Attempt to extract title and timestamp from the first NBT table, if it exists.
	firstTable := doc.Find("table.NBT").First()
	if firstTable.Length() > 0 {
		header := firstTable.Find("td")
		if header.Length() > 1 {
			if text := header.First().Text(); text != "" {
				scrape.Title = ptr(text)
			}

			if t, err := parseTimeToday(header.Last().Text()); err == nil {
				scrape.Timestamp = t
			}
		}
	}

	// Search the entire document for links.
	links := doc.Find("a")

	links.Each(func(_ int, s *goquery.Selection) {
		if attr, exists := s.Attr("href"); exists {
			if data := namedCapture(resultsLinkRegex, attr); data != nil {
				heat, err := strconv.Atoi(data["heat"])
				if err != nil {
					return // Skip silently if it fails to parse
				}

				round, err := strconv.Atoi(data["round"])
				if err != nil {
					return // Skip silently if it fails to parse
				}

				hr := models.HeatRound{Heat: heat, Round: round}

				switch data["type"] {
				case "p":
					if scrape.Practice == nil {
						scrape.Practice = newResultStatus()
					}
					scrape.Practice.Results = append(scrape.Practice.Results, hr)
				case "h":
					if scrape.Qualifying == nil {
						scrape.Qualifying = newResultStatus()
					}
					scrape.Qualifying.Results = append(scrape.Qualifying.Results, hr)
				case "f":
					if scrape.Finals == nil {
						scrape.Finals = newResultStatus()
					}
					scrape.Finals.Results = append(scrape.Finals.Results, hr)
				}
			} else if data := namedCapture(roundOverallsLinkRegex, attr); data != nil {
				round, err := strconv.Atoi(data["round"])
				if err != nil {
					return
				}

				switch data["type"] {
				case "p":
					if scrape.Practice == nil {
						scrape.Practice = newResultStatus()
					}
					if round > 0 {
						scrape.Practice.RoundOveralls = append(scrape.Practice.RoundOveralls, round)
					} else {
						scrape.Practice.Overall = true
					}
				case "h":
					if scrape.Qualifying == nil {
						scrape.Qualifying = newResultStatus()
					}
					if round > 0 {
						scrape.Qualifying.RoundOveralls = append(scrape.Qualifying.RoundOveralls, round)
					} else {
						scrape.Qualifying.Overall = true
					}
				case "f":
					if scrape.Finals == nil {
						scrape.Finals = newResultStatus()
					}
					if round > 0 {
						scrape.Finals.RoundOveralls = append(scrape.Finals.RoundOveralls, round)
					} else {
						scrape.Finals.Overall = true
					}
				}
			}
		}
	})

	if scrape.Practice != nil {
		sortHeatRounds(scrape.Practice.Results)
	}
	if scrape.Qualifying != nil {
		sortHeatRounds(scrape.Qualifying.Results)
	}
	if scrape.Finals != nil {
		sortHeatRounds(scrape.Finals.Results)
	}

	return scrape, nil
}
