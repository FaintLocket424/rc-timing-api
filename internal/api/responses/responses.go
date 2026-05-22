// Package responses contains methods that return the standard response
// format for the API
package responses

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// responseFormat represents the standard response format of the API.
type responseFormat struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

// RespondSuccess returns an OK 200 with data and a message.
func RespondSuccess(c *gin.Context, data any, message string) {
	c.JSON(http.StatusOK, responseFormat{
		Status:  "success",
		Data:    data,
		Message: message,
	})
}

// RespondAccepted returns an ACCEPTED 201 with a message.
func RespondAccepted(c *gin.Context, message string) {
	c.AbortWithStatusJSON(http.StatusAccepted, responseFormat{
		Status:  "success",
		Message: message,
	})
}

// RespondError returns an error code with a message.
func RespondError(c *gin.Context, code int, message string) {
	c.AbortWithStatusJSON(code, responseFormat{
		Status:  "error",
		Message: message,
	})
}
