package bbk

import "github.com/FaintLocket424/rc-timing-api/internal/models"

func (s BBKScraper) GetLiveTiming() *models.LiveTiming {
	rawHTML := ""
	return parseRaceResult(rawHTML)
}

func parseRaceResult(raw string) *models.LiveTiming {
	return &models.LiveTiming{}
}
