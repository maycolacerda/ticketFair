package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/maycolacerda/ticketfair/database"
	"github.com/maycolacerda/ticketfair/models"
)

// NewMerchant godoc
// @Summary Create a new merchant
// @Description Create a new merchant
// @Tags merchants
// @Accept json
// @Produce json
// @Param merchant body models.Merchant true "Merchant"
// @Success 201 {object} models.Merchant
// @Failure 400 {object} []string
// @Failure 500 {object} []string
// @Router /merchants/new/merchant [post]
func NewMerchant(c *gin.Context) {
	var merchant models.Merchant
	if err := c.ShouldBindJSON(&merchant); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := merchant.Validate(); len(err) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}
	if err := database.DB.Create(&merchant).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, merchant)
}
