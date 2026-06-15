package bbk

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/FaintLocket424/opengrid-bridge/internal/models"
	"github.com/PuerkitoBio/goquery"
)

func parseRaceResultsIndexHTML(body io.Reader) (*models.RaceResultsIndexScrape, error) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML body: %w", err)
	}

	scrape := &models.RaceResultsIndexScrape{}
	scrape.Title, scrape.Timestamp = parseTitleAndTimestamp(doc)

	// Dynamically locate the results index table based on contents.
	resultsTable := findResultsIndexTable(doc)
	if resultsTable == nil {
		return scrape, nil
	}

	// Limit link parsing exclusively to target links within the dynamic results table.
	links := resultsTable.Find("a")

	links.Each(func(_ int, s *goquery.Selection) {
		attr, exists := s.Attr("href")
		if !exists {
			return
		}

		if data := namedCapture(resultsLinkRegex, attr); data != nil {
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

// findResultsIndexTable searches for the table containing result grid indexes.
func findResultsIndexTable(doc *goquery.Document) *goquery.Selection {
	var targetTable *goquery.Selection
	doc.Find("table").Each(func(_ int, t *goquery.Selection) {
		if targetTable != nil {
			return
		}
		hasIndexIndicators := false
		t.Find("td").Each(func(_ int, td *goquery.Selection) {
			txt := strings.TrimSpace(td.Text())
			if txt == "Heats" || txt == "Practice" || txt == "Finals" || txt == "Overall" {
				hasIndexIndicators = true
			}
		})
		// Fallback detection: search inside table cells for links matching structural index regex patterns
		if !hasIndexIndicators {
			t.Find("a").Each(func(_ int, a *goquery.Selection) {
				if href, exists := a.Attr("href"); exists {
					if resultsLinkRegex.MatchString(href) || roundOverallsLinkRegex.MatchString(href) {
						hasIndexIndicators = true
					}
				}
			})
		}
		if hasIndexIndicators {
			targetTable = t
		}
	})
	return targetTable
}
