// tests/testutil/engine.go
package testutil

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RouterFromEngine converts a gin.Engine or any http.Handler to *gin.Engine
// for use with the fluent request builder.
func RouterFromEngine(e interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}) *gin.Engine {
	if g, ok := e.(*gin.Engine); ok {
		return g
	}
	// wrap in a new engine
	r := gin.New()
	r.Use(func(c *gin.Context) {
		e.ServeHTTP(c.Writer, c.Request)
		c.Abort()
	})
	return r
}
