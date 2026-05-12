package bbk

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/FaintLocket424/opengrid-bridge/internal/models"
)

func ptr[T any](v T) *T {
	return &v
}

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
		if err := res.Body.Close(); err != nil {
			slog.Error("Failed to close response body", "err", err, "url", url)
		}
		return nil, fmt.Errorf("server returned status: %d", res.StatusCode)
	}

	return res.Body, nil
}

func (s *BBKScraper) GetLiveTiming() (res *models.RaceResultScrape, err error) {
	body, fetchErr := s.fetchPage(s.Target + "/liveraceres.htm")
	if fetchErr != nil {
		return nil, fetchErr
	}
	defer func() {
		if closeErr := body.Close(); closeErr != nil {
			err = errors.Join(err, closeErr)
		}
	}()

	return parseRaceResult(body)
}

func (s *BBKScraper) GetPracticeRaceResult(practice, round int) (res *models.RaceResultScrape, err error) {
	if practice < 1 || round < 1 {
		return nil, fmt.Errorf("practice and round must be >= 1 (got %d, %d)", practice, round)
	}

	url := fmt.Sprintf("%s/p%dr%dres.htm", s.Target, practice, round)
	body, fetchErr := s.fetchPage(url)
	if fetchErr != nil {
		return nil, fetchErr
	}
	defer func() {
		if closeErr := body.Close(); closeErr != nil {
			err = errors.Join(err, closeErr)
		}
	}()

	return parseRaceResult(body)
}

func (s *BBKScraper) GetQualiRaceResult(heat, round int) (res *models.RaceResultScrape, err error) {
	if heat < 1 || round < 1 {
		return nil, fmt.Errorf("heat and round must be >= 1 (got %d, %d)", heat, round)
	}

	url := fmt.Sprintf("%s/h%dr%dres.htm", s.Target, heat, round)
	body, fetchErr := s.fetchPage(url)
	if fetchErr != nil {
		return nil, fetchErr
	}
	defer func() {
		if closeErr := body.Close(); closeErr != nil {
			err = errors.Join(err, closeErr)
		}
	}()

	return parseRaceResult(body)
}

func (s *BBKScraper) GetFinalRaceResult(final, leg int) (res *models.RaceResultScrape, err error) {
	if final < 1 || leg < 1 {
		return nil, fmt.Errorf("final and leg must be >= 1 (got %d, %d)", final, leg)
	}

	url := fmt.Sprintf("%s/f%dr%dres.htm", s.Target, final, leg)
	body, fetchErr := s.fetchPage(url)
	if fetchErr != nil {
		return nil, fetchErr
	}
	defer func() {
		if closeErr := body.Close(); closeErr != nil {
			err = errors.Join(err, closeErr)
		}
	}()

	return parseRaceResult(body)
}

func (s *BBKScraper) GetRaceResultsIndex() (res *models.RaceResultsIndexScrape, err error) {
	body, fetchErr := s.fetchPage(s.Target + "/liveresults.htm")
	if fetchErr != nil {
		return nil, fetchErr
	}
	defer func() {
		if closeErr := body.Close(); closeErr != nil {
			err = errors.Join(err, closeErr)
		}
	}()

	return parseRaceResultsIndex(body)
}
