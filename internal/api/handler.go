package api

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/FaintLocket424/opengrid-bridge/internal/api/dto"
	"github.com/FaintLocket424/opengrid-bridge/internal/storage"
	"github.com/gin-gonic/gin"
)

// Tracker represents a struct that can track a URL and populate a store
// with the data.
type Tracker interface {
	EnsureTracking(url string) bool
}

// Handler holds references to the current store and manager.
// It contains methods for handling different request types.
type Handler struct {
	store   storage.Store
	manager Tracker
}

// NewHandler creates a new Handler struct with references to a store and manager.
func NewHandler(store storage.Store, manager Tracker) *Handler {
	return &Handler{store, manager}
}

// GetLiveTiming handles HTTP requests to fetch live timing data.
// It ensures a tracking goroutine is running for the provided URL, then
// checks the cache. If tracking is initialising, it prompts the client to poll again.
func (h *Handler) GetLiveTiming(c *gin.Context) {
	url := c.Query("target_url")

	if url == "" {
		respondError(c, http.StatusBadRequest, "Missing required query parameter: target_url")
		return
	}

	newWorkerSpawned := h.manager.EnsureTracking(url)

	model, err := h.store.GetLiveTiming(url)
	if err != nil {
		if newWorkerSpawned {
			respondAccepted(c, "Starting tracking event, please poll again in a few seconds")
		} else {
			slog.Error("error getting live timing", "error", err)
			respondError(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondSuccess(c, dto.ToRaceResultDTO(model))
}

// GetPracticeRaceResult handles HTTP requests to fetch practice results.
// It ensures a tracking goroutine is running for the provided URL, then
// checks the cache. If tracking is initialising, it prompts the client to poll again.
func (h *Handler) GetPracticeRaceResult(c *gin.Context) {
	url := c.Query("target_url")

	if url == "" {
		respondError(c, http.StatusBadRequest, "Missing required query parameter: target_url")
		return
	}

	newWorkerSpawned := h.manager.EnsureTracking(url)

	heat, _ := strconv.Atoi(c.Param("heat"))
	round, _ := strconv.Atoi(c.Param("round"))

	model, err := h.store.GetPracticeRaceResult(url, heat, round)
	if err != nil {
		slog.Error("error getting practice race result", "heat", heat, "round", round, "error", err)
		if newWorkerSpawned {
			respondAccepted(c, "Starting tracking event, please poll again in a few seconds")
		} else {
			respondError(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondSuccess(c, dto.ToRaceResultDTO(model))
}

// GetQualiRaceResult handles HTTP requests to fetch practice results.
// It ensures a tracking goroutine is running for the provided URL, then
// checks the cache. If tracking is initialising, it prompts the client to poll again.
func (h *Handler) GetQualiRaceResult(c *gin.Context) {
	url := c.Query("target_url")

	if url == "" {
		respondError(c, http.StatusBadRequest, "Missing required query parameter: target_url")
		return
	}

	newWorkerSpawned := h.manager.EnsureTracking(url)

	heat, _ := strconv.Atoi(c.Param("heat"))
	round, _ := strconv.Atoi(c.Param("round"))

	model, err := h.store.GetQualiRaceResult(url, heat, round)
	if err != nil {
		slog.Error("error getting qualifying race result", "heat", heat, "round", round, "error", err)
		if newWorkerSpawned {
			respondAccepted(c, "Starting tracking event, please poll again in a few seconds")
		} else {
			respondError(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondSuccess(c, dto.ToRaceResultDTO(model))
}

// GetFinalRaceResult handles HTTP requests to fetch practice results.
// It ensures a tracking goroutine is running for the provided URL, then
// checks the cache. If tracking is initialising, it prompts the client to poll again.
func (h *Handler) GetFinalRaceResult(c *gin.Context) {
	url := c.Query("target_url")
	final, _ := strconv.Atoi(c.Param("final"))
	round, _ := strconv.Atoi(c.Param("round"))

	if url == "" {
		respondError(c, http.StatusBadRequest, "Missing required query parameter: target_url")
		return
	}

	newWorkerSpawned := h.manager.EnsureTracking(url)

	model, err := h.store.GetFinalRaceResult(url, final, round)
	if err != nil {
		slog.Error("error getting final race result", "final", final, "round", round, "error", err)
		if newWorkerSpawned {
			respondAccepted(c, "Starting tracking event, please poll again in a few seconds")
		} else {
			respondError(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondSuccess(c, dto.ToRaceResultDTO(model))
}
