// middlewares/merchant.go
package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/maycolacerda/ticketfair/services"
)

func MerchantMiddleware() gin.HandlerFunc {
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

		if claims.Role != services.RoleMerchant {
			abortForbidden(c)
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)
		c.Set("merchant_id", claims.MerchantID)
		c.Next()
	}
}
