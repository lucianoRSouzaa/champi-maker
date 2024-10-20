package web

import (
	"github.com/gin-gonic/gin"
)

func RespondWithJSON(c *gin.Context, code int, payload interface{}) {
	c.JSON(code, payload)
}

func RespondWithError(c *gin.Context, code int, message string) {
	RespondWithJSON(c, code, gin.H{"error": message})
}
