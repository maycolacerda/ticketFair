package middlewares

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/maycolacerda/ticketfair/services"
)

func ClientMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := services.ValidateToken(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized. Please Log in to use this feature."})
			c.Abort()
			slog.Warn("Unauthorized", "Path", c.Request.URL.Path)
			return
		}
		slog.Info("Authorized", "Method", c.Request.Method, "Path", c.Request.URL.Path)
		c.Next()
	}
}

func MerchantMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := services.ValidateToken(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized. Please Log in to use this feature."})
			c.Abort()
			slog.Warn("Unauthorized", "Path", c.Request.URL.Path)
			return
		}
		slog.Info("Authorized", "Method", c.Request.Method, "Path", c.Request.URL.Path)
		c.Next()
	}
}

func PublicMidleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		slog.Info("Public Access", c.Request.Method, c.Request.URL.Path)
		c.Next()
	}

}
