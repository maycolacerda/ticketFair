package controllers

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetHome godoc
//
//	@Summary		Show the home page.
//	@Description	get the home page.
//	@Tags			Home
//	@Accept			*/*
//	@Produce		json
//	@Success		200	{object}	map[string]string
//	@Router			/public/ [get]
func GetHome(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Welcome to the Ticket Fair API!"})

}

// HealthCheck godoc
//
//	@Summary		Show the status of the server.
//	@Description	get the status of the server.
//	@Tags			Health
//	@Accept			*/*
//	@Produce		json
//	@Success		200	{object}	map[string]interface{}
//	@Router			/public/health [get]
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"Server Status": "Ok"})
}

// NotFound godoc
//
//	@Summary		Handle not found routes.
//	@Description	Returns a 404 error for not found routes.
//	@Tags			NotFound
//	@Accept			*/*
//	@Produce		json
//	@Success		404	{object}	map[string]string
//	@Router			/not-found [get]
func NotFound(c *gin.Context) {
	slog.Warn("Route not found",
		"path", c.Request.URL.Path,
		"method", c.Request.Method,
	)
	c.JSON(http.StatusNotFound, gin.H{"error": "route not found"})
}
