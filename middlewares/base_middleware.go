// middlewares/base.go
package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/maycolacerda/ticketfair/services"
)

// extractBearer is shared across all middleware — do not redeclare in other files
func extractBearer(c *gin.Context) (string, error) {
	return services.ExtractBearerToken(c) // reuse the same logic from services/token.go
}

func abortUnauthorized(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
}

func abortForbidden(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
}
