// controllers/ticket_type.go
package controllers

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/maycolacerda/ticketfair/dto"
	"github.com/maycolacerda/ticketfair/services"
)

// ListTicketTypes godoc — public
//
//	@Summary		List ticket types for an event
//	@Description	Returns all active ticket types with availability and sale status
//	@Tags			Ticket Types
//	@Produce		json
//	@Param			id	path		string	true	"Event UUID"
//	@Success		200	{array}		dto.TicketTypeResponse
//	@Failure		404	{object}	map[string]string
//	@Router			/public/events/{id}/ticket-types [get]
func ListTicketTypes(c *gin.Context) {
	eventID := c.Param("id")

	types, err := services.GetTicketTypesByEvent(eventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch ticket types"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": types})
}

// ListAllTicketTypes godoc — merchant
//
//	@Summary		List all ticket types for merchant's event
//	@Description	Returns all ticket types including inactive ones
//	@Tags			Ticket Types
//	@Produce		json
//	@Param			id	path		string	true	"Event UUID"
//	@Success		200	{array}		dto.TicketTypeResponse
//	@Router			/merchant/events/{id}/ticket-types [get]
func ListAllTicketTypes(c *gin.Context) {
	merchantID, err := services.ExtractMerchantID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	eventID := c.Param("id")

	types, err := services.GetAllTicketTypesByEvent(merchantID, eventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch ticket types"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": types})
}

// CreateTicketType godoc
//
//	@Summary		Create a ticket type
//	@Description	Merchant — add a new ticket type to an event
//	@Tags			Ticket Types
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string						true	"Event UUID"
//	@Param			body	body		dto.CreateTicketTypeRequest	true	"Ticket type data"
//	@Success		201		{object}	dto.TicketTypeResponse
//	@Failure		400		{object}	map[string]string
//	@Failure		404		{object}	map[string]string
//	@Router			/merchant/events/{id}/ticket-types [post]
func CreateTicketType(c *gin.Context) {
	merchantID, err := services.ExtractMerchantID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	eventID := c.Param("id")

	var req dto.CreateTicketTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if err := validate.Struct(req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": formatValidationErrors(err)})
		return
	}

	tt, err := services.CreateTicketType(merchantID, eventID, req)
	if err != nil {
		slog.Warn("CreateTicketType failed",
			"merchant_id", merchantID,
			"event_id", eventID,
			"error", err.Error(),
		)
		switch {
		case errors.Is(err, services.ErrEventNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create ticket type"})
		}
		return
	}

	slog.Info("Ticket type created",
		"ticket_type_id", tt.TicketTypeID,
		"event_id", eventID,
	)
	c.JSON(http.StatusCreated, gin.H{"data": tt})
}

// UpdateTicketType godoc
//
//	@Summary		Update a ticket type
//	@Description	Merchant — update an existing ticket type
//	@Tags			Ticket Types
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string						true	"Event UUID"
//	@Param			ttid	path		string						true	"Ticket Type UUID"
//	@Param			body	body		dto.UpdateTicketTypeRequest	true	"Update data"
//	@Success		200		{object}	dto.TicketTypeResponse
//	@Failure		404		{object}	map[string]string
//	@Router			/merchant/events/{id}/ticket-types/{ttid} [put]
func UpdateTicketType(c *gin.Context) {
	merchantID, err := services.ExtractMerchantID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	eventID := c.Param("id")
	ticketTypeID := c.Param("ttid")

	var req dto.UpdateTicketTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if err := validate.Struct(req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": formatValidationErrors(err)})
		return
	}

	tt, err := services.UpdateTicketType(merchantID, eventID, ticketTypeID, req)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrTicketTypeNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, services.ErrNoFieldsToUpdate):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update ticket type"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": tt})
}

// DeleteTicketType godoc
//
//	@Summary		Delete a ticket type
//	@Description	Merchant — soft-delete a ticket type (only if no tickets sold)
//	@Tags			Ticket Types
//	@Produce		json
//	@Param			id		path		string	true	"Event UUID"
//	@Param			ttid	path		string	true	"Ticket Type UUID"
//	@Success		200		{object}	map[string]string
//	@Failure		400		{object}	map[string]string
//	@Failure		404		{object}	map[string]string
//	@Router			/merchant/events/{id}/ticket-types/{ttid} [delete]
func DeleteTicketType(c *gin.Context) {
	merchantID, err := services.ExtractMerchantID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	eventID := c.Param("id")
	ticketTypeID := c.Param("ttid")

	if err := services.DeleteTicketType(merchantID, eventID, ticketTypeID); err != nil {
		switch {
		case errors.Is(err, services.ErrTicketTypeNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, services.ErrTicketTypeHasSales):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete ticket type"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ticket type deleted successfully"})
}
