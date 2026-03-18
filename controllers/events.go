// controllers/events.go
package controllers

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/maycolacerda/ticketfair/dto"
	"github.com/maycolacerda/ticketfair/services"
)

// GetEvents godoc
//
//	@Summary		List all active events
//	@Description	Retrieve a paginated list of active upcoming events
//	@Tags			Events
//	@Produce		json
//	@Param			page	query		int	false	"Page number"	default(1)
//	@Param			limit	query		int	false	"Page size"		default(20)
//	@Success		200		{object}	dto.PaginatedEventsResponse
//	@Failure		500		{object}	map[string]string
//	@Router			/public/events [get]
func GetEvents(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	result, err := services.GetEvents(page, limit)
	if err != nil {
		slog.Error("Failed to fetch events", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch events"})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetEventByID godoc
//
//	@Summary		Get an event by ID
//	@Description	Retrieve a single event by its UUID
//	@Tags			Events
//	@Produce		json
//	@Param			id	path		string	true	"Event UUID"
//	@Success		200	{object}	dto.EventResponse
//	@Failure		404	{object}	map[string]string
//	@Router			/public/events/{id} [get]
func GetEventByID(c *gin.Context) {
	eventID := c.Param("id")

	event, err := services.GetEventByID(eventID)
	if err != nil {
		slog.Warn("Event not found", "event_id", eventID)
		c.JSON(http.StatusNotFound, gin.H{"error": "event not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": event})
}

// NewEvent godoc
//
//	@Summary		Create a new event
//	@Description	Create a new event under the authenticated merchant
//	@Tags			Events
//	@Accept			json
//	@Produce		json
//	@Param			event	body		dto.CreateEventRequest	true	"Event data"
//	@Success		201		{object}	dto.EventResponse
//	@Failure		400		{object}	map[string]string
//	@Failure		401		{object}	map[string]string
//	@Failure		422		{object}	map[string]interface{}
//	@Router			/merchant/events/new [post]
func NewEvent(c *gin.Context) {
	merchantID, err := services.ExtractTokenID(c)
	if err != nil {
		slog.Warn("Unauthorized event creation attempt")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req dto.CreateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Warn("Invalid request body", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := validate.Struct(req); err != nil {
		errs := formatValidationErrors(err)
		slog.Warn("Event validation failed", "errors", errs)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": errs})
		return
	}

	event, err := services.CreateEvent(merchantID, req)
	if err != nil {
		slog.Warn("Event creation failed", "merchant_id", merchantID, "error", err.Error())
		switch err.Error() {
		case "merchant not found":
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case "merchant account is disabled":
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		case "start_time must be in the future", "end_time must be after start_time":
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create event"})
		}
		return
	}

	slog.Info("Event created", "event_id", event.EventID, "merchant_id", merchantID)
	c.JSON(http.StatusCreated, gin.H{"data": event})
}

// UpdateEvent godoc
//
//	@Summary		Update an event
//	@Description	Update an event under the authenticated merchant
//	@Tags			Events
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string					true	"Event UUID"
//	@Param			event	body		dto.UpdateEventRequest	true	"Updated event data"
//	@Success		200		{object}	dto.EventResponse
//	@Failure		400		{object}	map[string]string
//	@Failure		401		{object}	map[string]string
//	@Failure		404		{object}	map[string]string
//	@Failure		422		{object}	map[string]interface{}
//	@Router			/merchant/events/{id} [put]
func UpdateEvent(c *gin.Context) {
	merchantID, err := services.ExtractTokenID(c)
	if err != nil {
		slog.Warn("Unauthorized event update attempt")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	eventID := c.Param("id")
	if eventID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing event id"})
		return
	}

	var req dto.UpdateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Warn("Invalid request body", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := validate.Struct(req); err != nil {
		errs := formatValidationErrors(err)
		slog.Warn("Event update validation failed", "errors", errs)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": errs})
		return
	}

	event, err := services.UpdateEvent(merchantID, eventID, req)
	if err != nil {
		slog.Warn("Event update failed", "event_id", eventID, "error", err.Error())
		switch err.Error() {
		case "event not found":
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case "no fields to update":
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case "start_time must be in the future", "end_time must be after start_time":
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update event"})
		}
		return
	}

	slog.Info("Event updated", "event_id", eventID, "merchant_id", merchantID)
	c.JSON(http.StatusOK, gin.H{"data": event})
}
