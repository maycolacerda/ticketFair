package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetHome(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Welcome to the Ticket Fair API!"})
}

func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "Tá saudável"})
}

func NotFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{"error": "Página não encontrada"})
}
