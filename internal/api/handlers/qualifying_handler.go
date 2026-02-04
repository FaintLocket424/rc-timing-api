package handlers

import (
	"net/http"
	"strconv"

	"github.com/FaintLocket424/rc-timing-api/internal/api/middleware"
	"github.com/FaintLocket424/rc-timing-api/internal/service"
	"github.com/gin-gonic/gin"
)

type QualifyingHandler struct {
	Store *service.DataStore
}

func NewQualifyingHandler(s *service.DataStore) *QualifyingHandler {
	return &QualifyingHandler{Store: s}
}

func (h *QualifyingHandler) GetHeatRoundResult(c *gin.Context) {
	url := middleware.GetURL(c)
	scrp := middleware.GetScraper(c)

	heatID, _ := strconv.Atoi(c.Param("id"))
	roundID, _ := strconv.Atoi(c.Param("round_id"))

	raw, err := h.Store.GetHeatResult(url, scrp, heatID, roundID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	dto := service.MapRawToDTO(raw)

	c.JSON(http.StatusOK, dto)
}
