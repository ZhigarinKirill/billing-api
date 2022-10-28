package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// Error - error structure
type Error struct {
	ErrorMessage string `json:"error_message"`
}

// newErrorResponse - error wrapper
func newErrorResponse(c *gin.Context, statusCode int, message string) {
	log.Error().Msg(message)
	c.AbortWithStatusJSON(statusCode, &Error{ErrorMessage: message})
}
