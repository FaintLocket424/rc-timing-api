package api

import (
	"net/http"

	"github.com/FaintLocket424/rc-timing-api/internal/storage"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	store storage.Store
}

func NewHandler(store storage.Store) *Handler {
	return &Handler{store}
}

func (h *Handler) GetLiveTiming(c *gin.Context) {
	url := c.Query("target_url")

	if url == "" {
		c.AbortWithStatus(http.StatusBadRequest)
	}

	model, err := h.store.GetLiveTiming(url)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	}

	c.JSON(http.StatusOK, model)
}
