package handlers

import (
	"net/http"

	"github.com/FaintLocket424/rc-timing-api/internal/api/middleware"
	"github.com/FaintLocket424/rc-timing-api/internal/service"
	"github.com/gin-gonic/gin"
)

type EventHandler struct {
	Store *service.DataStore
}

func (h *EventHandler) GetMeta(c *gin.Context) {
	url := middleware.GetURL(c)
	scraper := middleware.GetScraper(c)

	raw, err := h.Store.GetEventMeta(url, scraper)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	dto := service.MapCachedMetaToDTO(raw)

	c.JSON(http.StatusOK, dto)
}

func NewEventHandler(s *service.DataStore) *EventHandler {
	return &EventHandler{Store: s}
}
