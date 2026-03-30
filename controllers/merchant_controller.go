// controllers/merchant.go
package controllers

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/maycolacerda/ticketfair/dto"
	"github.com/maycolacerda/ticketfair/services"
)

// NewMerchant godoc
//
//	@Summary		Create a new merchant
//	@Description	Create a new merchant account
//	@Tags			Merchants
//	@Accept			json
//	@Produce		json
//	@Param			merchant	body		dto.CreateMerchantRequest	true	"Merchant data"
//	@Success		201			{object}	dto.MerchantResponse
//	@Failure		400			{object}	map[string]string
//	@Failure		409			{object}	map[string]string
//	@Failure		422			{object}	map[string]interface{}
//	@Router			/public/merchant/register [post]
func NewMerchant(c *gin.Context) {
	var req dto.CreateMerchantRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Warn("Invalid request body", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := validate.Struct(req); err != nil {
		errs := formatValidationErrors(err)
		slog.Warn("Merchant validation failed", "errors", errs)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": errs})
		return
	}

	merchant, err := services.CreateMerchant(req)
	if err != nil {
		slog.Warn("Merchant creation failed", "error", err.Error())
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	slog.Info("Merchant created", "merchant_id", merchant.MerchantID)
	c.JSON(http.StatusCreated, gin.H{"data": merchant})
}

// UpdateMerchant godoc
//
//	@Summary		Update a merchant
//	@Description	Update the authenticated merchant's details
//	@Tags			Merchants
//	@Accept			json
//	@Produce		json
//	@Param			merchant	body		dto.UpdateMerchantRequest	true	"Updated merchant data"
//	@Success		200			{object}	dto.MerchantResponse
//	@Failure		400			{object}	map[string]string
//	@Failure		401			{object}	map[string]string
//	@Failure		404			{object}	map[string]string
//	@Failure		422			{object}	map[string]interface{}
//	@Router			/merchant/update [put]
func UpdateMerchant(c *gin.Context) {
	// ID comes from JWT — never trust client-supplied IDs for ownership
	merchantID, err := services.ExtractTokenID(c)
	if err != nil {
		slog.Warn("Unauthorized update attempt")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req dto.UpdateMerchantRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Warn("Invalid request body", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := validate.Struct(req); err != nil {
		errs := formatValidationErrors(err)
		slog.Warn("Merchant update validation failed", "errors", errs)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": errs})
		return
	}

	merchant, err := services.UpdateMerchant(merchantID, req)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrMerchantNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, services.ErrNoFieldsToUpdate):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update merchant"})
		}
		return
	}
	slog.Info("Merchant updated", "merchant_id", merchantID)
	c.JSON(http.StatusOK, gin.H{"data": merchant})
}
