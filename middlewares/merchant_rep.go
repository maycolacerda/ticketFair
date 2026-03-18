// middlewares/merchant_rep.go
package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/maycolacerda/ticketfair/services"
)

func MerchantRepMiddleware(allowedRoles ...string) gin.HandlerFunc {
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

		allowed := false
		for _, role := range allowedRoles {
			if claims.Role == role {
				allowed = true
				break
			}
		}
		if !allowed {
			abortForbidden(c)
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)
		c.Set("merchant_id", claims.MerchantID)
		c.Next()
	}
}
