package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/FaintLocket424/rc-timing-api/internal/api/middleware"
	_ "github.com/FaintLocket424/rc-timing-api/internal/models"
	"github.com/FaintLocket424/rc-timing-api/internal/service"
	"github.com/gin-gonic/gin"
)

type PracticeHandler struct {
	Store *service.DataStore
}

func NewPracticeHandler(s *service.DataStore) *PracticeHandler {
	return &PracticeHandler{Store: s}
}

// GetHeatRoundResult godoc
// @Summary      Get practice heat results
// @Description  Retrieves results for a specific practice heat and round using the target timing software.
// @Tags         practice
// @Produce      json
// @Param        url       query     string  true  "Target Timing URL (URL-encoded)"
// @Param        software  query     string  false "Timing software type (default: bbk)"
// @Param        id        path      int     true  "Heat Number"
// @Param        round_id  path      int     true  "Round Number"
// @Success      200  {object}  models.HeatResultDTO
// @Failure      400  {object}  map[string]string "Invalid Heat or Round ID"
// @Failure      500  {object}  map[string]string "Scraper or Server Error"
// @Router       /event/practice/heat/{id}/round/{round_id} [get]
func (h *PracticeHandler) GetHeatRoundResult(c *gin.Context) {
	url := middleware.GetURL(c)
	scraper := middleware.GetScraper(c)

	// Retrieve the Heat number from the URL and convert it to an integer
	heatIDStr := c.Param("id")
	heatID, err := strconv.Atoi(heatIDStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid id: %s", heatIDStr)})
		return
	}

	// Retrieve the Round number from the URL and convert it to an integer
	roundIDStr := c.Param("round_id")
	roundID, err := strconv.Atoi(roundIDStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid round id: %s", roundIDStr)})
		return
	}

	// Retrieve the requested heat from the data store.
	raw, err := h.Store.GetPracticeHeatResult(url, scraper, heatID, roundID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Map the raw data to a Data Transfer Object to be sent back.
	dto := service.MapCachedHeatResultToDTO(raw)

	c.JSON(http.StatusOK, dto)
}
