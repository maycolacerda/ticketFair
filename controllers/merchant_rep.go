// controllers/merchantRep.go
package controllers

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/maycolacerda/ticketfair/dto"
	"github.com/maycolacerda/ticketfair/services"
)

// NewMerchantRep godoc
//
//	@Summary		Create a new merchant representative
//	@Description	Create a new rep under the authenticated merchant
//	@Tags			Merchant Representatives
//	@Accept			json
//	@Produce		json
//	@Param			merchantRep	body		dto.CreateMerchantRepRequest	true	"Merchant rep data"
//	@Success		201			{object}	dto.MerchantRepResponse
//	@Failure		400			{object}	map[string]string
//	@Failure		401			{object}	map[string]string
//	@Failure		404			{object}	map[string]string
//	@Failure		409			{object}	map[string]string
//	@Failure		422			{object}	map[string]interface{}
//	@Router			/merchant/rep/new [post]
func NewMerchantRep(c *gin.Context) {
	// MerchantID from JWT — rep is scoped to the authenticated merchant
	merchantID, err := services.ExtractTokenID(c)
	if err != nil {
		slog.Warn("Unauthorized merchant rep creation attempt")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req dto.CreateMerchantRepRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Warn("Invalid request body", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := validate.Struct(req); err != nil {
		errs := formatValidationErrors(err)
		slog.Warn("Merchant rep validation failed", "errors", errs)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": errs})
		return
	}

	rep, err := services.CreateMerchantRep(merchantID, req)
	if err != nil {
		slog.Warn("Merchant rep creation failed", "merchant_id", merchantID, "error", err.Error())
		switch err.Error() {
		case "merchant not found":
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case "email already in use for this merchant":
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create merchant representative"})
		}
		return
	}

	slog.Info("Merchant rep created", "merchant_rep_id", rep.MerchantRepID, "merchant_id", merchantID)
	c.JSON(http.StatusCreated, gin.H{"data": rep})
}

// UpdateMerchantRep godoc
//
//	@Summary		Update a merchant representative
//	@Description	Update a rep's details under the authenticated merchant
//	@Tags			Merchant Representatives
//	@Accept			json
//	@Produce		json
//	@Param			id			path		string							true	"Merchant Rep ID"
//	@Param			merchantRep	body		dto.UpdateMerchantRepRequest	true	"Updated rep data"
//	@Success		200			{object}	dto.MerchantRepResponse
//	@Failure		400			{object}	map[string]string
//	@Failure		401			{object}	map[string]string
//	@Failure		404			{object}	map[string]string
//	@Failure		422			{object}	map[string]interface{}
//	@Router			/merchant/rep/{id} [put]
func UpdateMerchantRep(c *gin.Context) {
	merchantID, err := services.ExtractTokenID(c)
	if err != nil {
		slog.Warn("Unauthorized merchant rep update attempt")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	repID := c.Param("id")
	if repID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing rep id"})
		return
	}

	var req dto.UpdateMerchantRepRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Warn("Invalid request body", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := validate.Struct(req); err != nil {
		errs := formatValidationErrors(err)
		slog.Warn("Merchant rep update validation failed", "errors", errs)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": errs})
		return
	}

	rep, err := services.UpdateMerchantRep(merchantID, repID, req)
	if err != nil {
		slog.Warn("Merchant rep update failed", "rep_id", repID, "error", err.Error())
		switch err.Error() {
		case "merchant representative not found":
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case "no fields to update":
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update merchant representative"})
		}
		return
	}

	slog.Info("Merchant rep updated", "merchant_rep_id", repID, "merchant_id", merchantID)
	c.JSON(http.StatusOK, gin.H{"data": rep})
}
