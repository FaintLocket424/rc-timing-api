package api

import (
	"net/http"

	"github.com/FaintLocket424/rc-timing-api/internal/manager"
	"github.com/FaintLocket424/rc-timing-api/internal/storage"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	store   storage.Store
	manager *manager.Manager
}

func NewHandler(store storage.Store, manager *manager.Manager) *Handler {
	return &Handler{store, manager}
}

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
