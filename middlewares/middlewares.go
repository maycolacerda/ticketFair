package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/maycolacerda/ticketfair/services"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := services.ValidateToken(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized. Please Log in to use this feature."})
			c.Abort()
			return
		}
		c.Next()
	}
}
