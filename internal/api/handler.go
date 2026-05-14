package api

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/FaintLocket424/opengrid-bridge/internal/manager"
	"github.com/FaintLocket424/opengrid-bridge/internal/storage"
	"github.com/gin-gonic/gin"
)

// Handler holds referencefs to the current store and manager.
// It contains methods for handling different request types.
type Handler struct {
	store   storage.Store
	manager *manager.Manager
}

// NewHandler creates a new Handler struct with references to a store and manager.
func NewHandler(store storage.Store, manager *manager.Manager) *Handler {
	return &Handler{store, manager}
}

// GetLiveTiming handles HTTP requests to fetch live timing data.
// It checks the store for cached data and ensures a tracking goroutine is running
// for the provided URL. If tracking is initializing, it prompts the client to poll again.
func (h *Handler) GetLiveTiming(c *gin.Context) {
	url := c.Query("target_url")

	if url == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	newWorkerSpawned := h.manager.EnsureTracking(url)

	model, err := h.store.GetLiveTiming(url)
	if err != nil {
		if newWorkerSpawned {
			c.AbortWithStatusJSON(http.StatusAccepted, gin.H{"message": "Starting tracking event, please poll again in a few seconds"})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		}
		return
	}

	c.JSON(http.StatusOK, model)
}

func (h *Handler) GetPracticeRaceResult(c *gin.Context) {
	url := c.Query("target_url")
	heat, _ := strconv.Atoi(c.Param("heat"))
	round, _ := strconv.Atoi(c.Param("round"))

	if url == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	newWorkerSpawned := h.manager.EnsureTracking(url)

	model, err := h.store.GetPracticeRaceResult(url, heat, round)
	if err != nil {
		slog.Error("error getting practice race result", "error", err)
		if newWorkerSpawned {
			c.AbortWithStatusJSON(http.StatusAccepted, gin.H{"message": "Starting tracking event, please poll again in a few seconds"})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, model)
}

func (h *Handler) GetQualiRaceResult(c *gin.Context) {
	url := c.Query("target_url")
	heat, _ := strconv.Atoi(c.Param("heat"))
	round, _ := strconv.Atoi(c.Param("round"))

	if url == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	newWorkerSpawned := h.manager.EnsureTracking(url)

	model, err := h.store.GetQualiRaceResult(url, heat, round)
	if err != nil {
		slog.Error("error getting quali race result", "error", err)
		if newWorkerSpawned {
			c.AbortWithStatusJSON(http.StatusAccepted, gin.H{"message": "Starting tracking event, please poll again in a few seconds"})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, model)
}

func (h *Handler) GetFinalRaceResult(c *gin.Context) {
	url := c.Query("target_url")
	heat, _ := strconv.Atoi(c.Param("heat"))
	round, _ := strconv.Atoi(c.Param("round"))

	if url == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	newWorkerSpawned := h.manager.EnsureTracking(url)

	model, err := h.store.GetFinalRaceResult(url, heat, round)
	if err != nil {
		slog.Error("error getting final race result", "error", err)
		if newWorkerSpawned {
			c.AbortWithStatusJSON(http.StatusAccepted, gin.H{"message": "Starting tracking event, please poll again in a few seconds"})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, model)
}
