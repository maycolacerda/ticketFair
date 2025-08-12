package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/maycolacerda/ticketfair/database"
	"github.com/maycolacerda/ticketfair/models"
)

// NewMerchantRep godoc
// @Summary		Create a new merchant representative.
// @Description	Create a new merchant representative with the provided details.
// @Tags			Merchant Representatives
// @Accept			json
// @Produce		json
// @Param			merchantRep	body	models.MerchantRep	true	"Merchant Representative data"
// @Success		201	{object}	models.MerchantRep
// @Failure		400	{object}	map[string]string
// @Failure		500	{object}	map[string]string
// @Router			/merchant/new/rep [post]
func NewMerchantRep(c *gin.Context) {
	var merchantRep models.MerchantRep
	if err := c.ShouldBindJSON(&merchantRep); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := merchantRep.Validate(); len(err) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}
	if err := database.DB.Create(&merchantRep).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, merchantRep)
}
