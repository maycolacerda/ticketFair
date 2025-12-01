package controllers

import (
	"log/slog"
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
		slog.Warn("Invalid request body", "details", err.Error)
		return
	}
	if err := merchant.Validate(); len(err) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		slog.Warn("Invalid request body", "details", err)
		return
	}
	if err := database.DB.Create(&merchant).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		slog.Error("Failed to create merchant", "details", err.Error)
		return
	}
	c.JSON(http.StatusCreated, merchant)
	slog.Info("Merchant created", "merchant_id", merchant.MerchantID)
}

// UpdateMerchant godoc
// @Summary Update a merchant
// @Description Update a merchant
// @Tags merchants
// @Accept json
// @Produce json
// @Param merchant body models.Merchant true "Merchant"
// @Success 200 {object} models.Merchant
// @Failure 400
// @Failure 500
// @router /merchant/update [post]
func UpdateMerchant(c *gin.Context) {
	var merchant models.Merchant
	if err := c.ShouldBindJSON(&merchant); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		slog.Warn("Invalid request body", "details", err.Error)
		return
	}
	if err := merchant.Validate(); len(err) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		slog.Warn("Invalid request body", "details", err)
		return
	}
	if err := database.DB.Save(&merchant).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		slog.Error("Failed to update merchant", "details", err.Error)
		return
	}
	c.JSON(http.StatusOK, merchant)
	slog.Info("Merchant updated", "merchant_id", merchant.MerchantID)

}
