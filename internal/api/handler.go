package api

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/FaintLocket424/opengrid-bridge/internal/api/dto"
	"github.com/FaintLocket424/opengrid-bridge/internal/api/responses"
	"github.com/FaintLocket424/opengrid-bridge/internal/storage"
	"github.com/gin-gonic/gin"
)

// Tracker represents a struct that can track a URL and populate a store
// with the data.
type Tracker interface {
	EnsureTracking(url string) bool
}

// Handler holds references to the current store and tracker.
// It contains methods for handling different request types.
type Handler struct {
	store   storage.Store
	tracker Tracker
}

// NewHandler creates a new Handler struct with references to a store and tracker.
func NewHandler(store storage.Store, tracker Tracker) *Handler {
	return &Handler{store, tracker}
}

// getTargetURL extracts the key set by the ExtractTargetURL middleware.
func getTargetURL(c *gin.Context) string {
	if val, ok := c.Get("target_url"); ok {
		return val.(string)
	}
	return ""
}

// GetLiveTiming handles HTTP requests to fetch live timing data.
// It ensures a tracking goroutine is running for the provided URL, then
// checks the cache. If tracking is initialising, it prompts the client to poll again.
func (h *Handler) GetLiveTiming(c *gin.Context) {
	url := getTargetURL(c)

	if h.tracker.EnsureTracking(url) {
		responses.RespondAccepted(c, "Starting tracking event, please poll again in a few seconds")
	} else {
		if model, err := h.store.GetLiveTiming(url); err == nil {
			responses.RespondSuccess(c, dto.ToRaceResultDTO(model), "")
		} else {
			slog.Error("error getting live timing", "error", err)
			responses.RespondError(c, http.StatusInternalServerError, err.Error())
		}
	}
}

// GetPracticeRaceResult handles HTTP requests to fetch practice results.
// It ensures a tracking goroutine is running for the provided URL, then
// checks the cache. If tracking is initialising, it prompts the client to poll again.
func (h *Handler) GetPracticeRaceResult(c *gin.Context) {
	url := getTargetURL(c)

	heat, _ := strconv.Atoi(c.Param("heat"))
	round, _ := strconv.Atoi(c.Param("round"))

	if h.tracker.EnsureTracking(url) {
		responses.RespondAccepted(c, "Starting tracking event, please poll again in a few seconds")
	} else {
		if model, err := h.store.GetPracticeRaceResult(url, heat, round); err == nil {
			responses.RespondSuccess(c, dto.ToRaceResultDTO(model), "")
		} else {
			slog.Error("error getting practice race result", "error", err)
			responses.RespondError(c, http.StatusInternalServerError, err.Error())
		}
	}
}

// GetQualiRaceResult handles HTTP requests to fetch practice results.
// It ensures a tracking goroutine is running for the provided URL, then
// checks the cache. If tracking is initialising, it prompts the client to poll again.
func (h *Handler) GetQualiRaceResult(c *gin.Context) {
	url := getTargetURL(c)

	heat, _ := strconv.Atoi(c.Param("heat"))
	round, _ := strconv.Atoi(c.Param("round"))

	if h.tracker.EnsureTracking(url) {
		responses.RespondAccepted(c, "Starting tracking event, please poll again in a few seconds")
	} else {
		if model, err := h.store.GetQualiRaceResult(url, heat, round); err == nil {
			responses.RespondSuccess(c, dto.ToRaceResultDTO(model), "")
		} else {
			slog.Error("error getting qualifying race result", "error", err)
			responses.RespondError(c, http.StatusInternalServerError, err.Error())
		}
	}
}

// GetFinalRaceResult handles HTTP requests to fetch practice results.
// It ensures a tracking goroutine is running for the provided URL, then
// checks the cache. If tracking is initialising, it prompts the client to poll again.
func (h *Handler) GetFinalRaceResult(c *gin.Context) {
	url := getTargetURL(c)
	final, _ := strconv.Atoi(c.Param("final"))
	round, _ := strconv.Atoi(c.Param("round"))

	if h.tracker.EnsureTracking(url) {
		responses.RespondAccepted(c, "Starting tracking event, please poll again in a few seconds")
	} else {
		if model, err := h.store.GetFinalRaceResult(url, final, round); err == nil {
			responses.RespondSuccess(c, dto.ToRaceResultDTO(model), "")
		} else {
			slog.Error("error getting final race result", "error", err)
			responses.RespondError(c, http.StatusInternalServerError, err.Error())
		}
	}
}
