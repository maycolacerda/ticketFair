package controllers

import (
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
//	@Router			/ [get]
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
//	@Router			/health [get]
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "Tá saudável"})
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
	c.JSON(http.StatusNotFound, gin.H{"error": "Página não encontrada"})
}
