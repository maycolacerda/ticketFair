// controllers/image.go
package controllers

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/maycolacerda/ticketfair/services"
)

// UploadEventImage godoc
//
//	@Summary		Upload event image
//	@Description	Merchant — upload a cover image for an event (JPEG, PNG or WebP, max 5MB)
//	@Tags			Events
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			id		path		string	true	"Event UUID"
//	@Param			image	formData	file	true	"Image file"
//	@Success		200		{object}	map[string]string
//	@Failure		400		{object}	map[string]string
//	@Failure		401		{object}	map[string]string
//	@Failure		404		{object}	map[string]string
//	@Failure		413		{object}	map[string]string
//	@Router			/merchant/events/{id}/image [post]
func UploadEventImage(c *gin.Context) {
	merchantID, err := services.ExtractMerchantID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	eventID := c.Param("id")

	// Parse multipart file — 5MB limit enforced at service level too
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "image file is required"})
		return
	}
	defer file.Close()

	// Upload to S3 and get URL
	imageURL, err := services.UploadEventImage(file, header)
	if err != nil {
		slog.Warn("Image upload failed",
			"event_id", eventID,
			"merchant_id", merchantID,
			"error", err.Error(),
		)
		switch {
		case errors.Is(err, services.ErrImageTooLarge):
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": err.Error()})
		case errors.Is(err, services.ErrInvalidImageType):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, services.ErrS3NotConfigured):
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "image storage unavailable"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upload image"})
		}
		return
	}

	// Persist URL on the event (scoped to merchant)
	event, err := services.SetEventImage(merchantID, eventID, imageURL)
	if err != nil {
		// Upload succeeded but DB write failed — best effort cleanup
		_ = services.DeleteEventImage(imageURL)
		switch {
		case errors.Is(err, services.ErrEventNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save image"})
		}
		return
	}

	slog.Info("Event image uploaded",
		"event_id", eventID,
		"merchant_id", merchantID,
		"image_url", imageURL,
	)

	c.JSON(http.StatusOK, gin.H{
		"message":   "image uploaded successfully",
		"image_url": event.ImageURL,
	})
}

// DeleteEventImage godoc
//
//	@Summary		Delete event image
//	@Description	Merchant — remove the cover image from an event
//	@Tags			Events
//	@Produce		json
//	@Param			id	path		string	true	"Event UUID"
//	@Success		200	{object}	map[string]string
//	@Failure		401	{object}	map[string]string
//	@Failure		404	{object}	map[string]string
//	@Router			/merchant/events/{id}/image [delete]
func DeleteEventImageHandler(c *gin.Context) {
	merchantID, err := services.ExtractMerchantID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	eventID := c.Param("id")

	if err := services.RemoveEventImage(merchantID, eventID); err != nil {
		switch {
		case errors.Is(err, services.ErrEventNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete image"})
		}
		return
	}

	slog.Info("Event image deleted", "event_id", eventID, "merchant_id", merchantID)
	c.JSON(http.StatusOK, gin.H{"message": "image deleted successfully"})
}
