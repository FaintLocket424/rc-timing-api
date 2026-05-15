package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ResponseFormat represents the standard response format of the API.
type ResponseFormat struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func respondSuccess(c *gin.Context, data any) {
	c.JSON(http.StatusOK, ResponseFormat{
		Status: "success",
		Data:   data,
	})
}

func respondAccepted(c *gin.Context, message string) {
	c.AbortWithStatusJSON(http.StatusAccepted, ResponseFormat{
		Status:  "success",
		Message: message,
	})
}

func respondError(c *gin.Context, code int, message string) {
	c.AbortWithStatusJSON(code, ResponseFormat{
		Status:  "error",
		Message: message,
	})
}
