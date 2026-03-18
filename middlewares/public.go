package middlewares

import (
	"log/slog"

	"github.com/gin-gonic/gin"
)

func PublicMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		slog.Info("Public Access", c.Request.Method, c.Request.URL.Path)
		c.Next()
	}

}
