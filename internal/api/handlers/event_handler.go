package handlers

import (
	"net/http"

	"github.com/FaintLocket424/rc-timing-api/internal/scraper"
	"github.com/FaintLocket424/rc-timing-api/internal/service"
	"github.com/gin-gonic/gin"

	_ "github.com/FaintLocket424/rc-timing-api/internal/models"
)

type EventHandler struct {
	Store *service.DataStore
}

func NewEventHandler(store *service.DataStore) *EventHandler {
	return &EventHandler{Store: store}
}

// GetLiveRace godoc
// @Summary      Get live race
// @Description  Returns the currently active heat for the provided URL
// @Tags         racing
// @Produce      json
// @Param        url query string true "Target Timing URL"
// @Success      200 {object} models.Heat
// @Router       /live [get]
func (h *EventHandler) GetLiveRace(c *gin.Context) {
	targetURL := c.Query("url")
	event, err := h.Store.GetEvent(targetURL, scraper.Scrape)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, race := range event.Heats {
		if race.Status == "Live" {
			c.JSON(http.StatusOK, race)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "No live race found"})
}

func GetHeatByType(store *service.DataStore) func(c *gin.Context) {
	return func(c *gin.Context) {
		data, _ := store.GetHeatsByType(c.Query("url"), c.Param("type"), scraper.Scrape)
		c.JSON(http.StatusOK, data)
	}
}

func GetSchedule(store *service.DataStore) func(c *gin.Context) {
	return func(c *gin.Context) {
		data, err := store.GetSchedule(c.Query("url"), scraper.Scrape)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Set("last_updated", data[0].ScheduledAt.Unix())
		c.JSON(http.StatusOK, data)
	}
}
