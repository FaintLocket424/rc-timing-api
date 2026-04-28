package bbk

import (
	"fmt"
	"io"
	"net/http"

	"github.com/FaintLocket424/rc-timing-api/internal/models"
)

type HTTPClient interface {
	Get(url string) (*http.Response, error)
}

type BBKScraper struct {
	Target string
	Client HTTPClient
}

func (s *BBKScraper) fetchPage(url string) (io.ReadCloser, error) {
	res, err := s.Client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("network error: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		res.Body.Close()
		return nil, fmt.Errorf("server returned status: %d", res.StatusCode)
	}

	return res.Body, nil
}

func (s *BBKScraper) GetLiveTiming() (*models.LiveTimingScrape, error) {
	body, err := s.fetchPage(s.Target + "/liveraceres.htm")
	if err != nil {
		return nil, err
	}
	defer body.Close()

	return parseRaceResult(body)
}

func (s *BBKScraper) GetQualiHeatResult(heat, round int) (*models.LiveTimingScrape, error) {
	if heat < 1 || round < 1 {
		return nil, fmt.Errorf("heat and round must be >= 1 (got %d, %d)", heat, round)
	}

	url := fmt.Sprintf("%s/h%dr%dres.htm", s.Target, heat, round)
	body, err := s.fetchPage(url)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	return parseRaceResult(body)
}
