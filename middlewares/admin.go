// middlewares/admin.go
package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/maycolacerda/ticketfair/services"
)

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := extractBearer(c)
		if err != nil {
			abortUnauthorized(c)
			return
		}

		claims, err := services.ParseToken(tokenStr)
		if err != nil {
			abortUnauthorized(c)
			return
		}

		if claims.Role != services.RoleAdmin {
			abortForbidden(c)
			return
		}

		c.Set("admin_id", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}
