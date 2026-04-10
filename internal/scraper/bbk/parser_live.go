package bbk

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/FaintLocket424/rc-timing-api/internal/models"
)

func (s BBKScraper) GetLiveTiming() (*models.LiveTiming, error) {
	res, err := http.Get(s.Target + "/liveraceres.htm")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code error: %d, %s", res.StatusCode, res.Status)
	}

	return parseRaceResult(res.Body)
}

func parseRaceResult(body io.Reader) (*models.LiveTiming, error) {
	return &models.LiveTiming{}, errors.New("Not implemented")
}
