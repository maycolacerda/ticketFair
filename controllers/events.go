package controllers

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/maycolacerda/ticketfair/database"
	"github.com/maycolacerda/ticketfair/models"
)

// NewEvent godoc
//
//		@Summary		Create a new event.
//		@Description	Create a new event with the provided details.
//		@Tags			Events
//		@Accept			json
//		@Produce		json
//		@Param			event	body	models.Event	true	"Event data"
//		@Success		200	{object}	map[string]string
//		@Failure		400	{object}	map[string]string
//	 	@Failure 500 {object}  []string
//		@Router /merchant/events/event/new [post]
func NewEvent(c *gin.Context) {
	var event models.Event
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body", "details": err.Error()})
		slog.Warn("Invalid request body", "details", err.Error)
		return
	}
	if err := event.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		slog.Warn("Invalid request body", "details", err)
		return
	}
	if err := database.DB.Create(&event).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create event"})
		slog.Error("Failed to create event", "details", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Event created successfully"})
	slog.Info("Event created", "event_id", event.EventID)
}

// NewEvent godoc
//
//		@Summary		Create a new event.
//		@Description	Create a new event with the provided details.
//		@Tags			Events
//		@Accept			json
//		@Produce		json
//		@Param			event	body	models.Event	true	"Event data"
//		@Success		200	{object}	map[string]string
//		@Failure		400	{object}	map[string]string
//	 	@Failure 500 {object}  []string
//		@Router /merchant/events/event/update[post]
func UpdateEvent(c *gin.Context) {
	var event models.Event
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body", "details": err.Error()})
		slog.Warn("Invalid request body", "details", err.Error)
		return
	}
	if err := event.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		slog.Warn("Invalid request body", "details", err)
		return
	}
	if err := database.DB.Save(&event).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update event"})
		slog.Error("Failed to update event", "details", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Event updated successfully"})
	slog.Info("Event updated", "event_id", event.EventID)

}
