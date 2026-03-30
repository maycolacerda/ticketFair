// controllers/ticket.go
package controllers

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/maycolacerda/ticketfair/services"
)

// GetMyTickets godoc
//
//	@Summary		List user tickets
//	@Description	Retrieve paginated ticket list for the authenticated user
//	@Tags			Tickets
//	@Produce		json
//	@Param			page	query		int	false	"Page number"	default(1)
//	@Param			limit	query		int	false	"Page size"		default(20)
//	@Success		200		{object}	dto.PaginatedTicketsResponse
//	@Failure		401		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/private/tickets [get]
func GetMyTickets(c *gin.Context) {
	userID, err := services.ExtractTokenID(c)
	if err != nil {
		slog.Warn("Unauthorized ticket list attempt")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	result, err := services.GetUserTickets(userID, page, limit)
	if err != nil {
		slog.Error("Failed to fetch tickets", "user_id", userID, "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch tickets"})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetTicketByID godoc
//
//	@Summary		Get a ticket by ID
//	@Description	Retrieve a single ticket by UUID for the authenticated user
//	@Tags			Tickets
//	@Produce		json
//	@Param			id	path		string	true	"Ticket UUID"
//	@Success		200	{object}	dto.TicketResponse
//	@Failure		401	{object}	map[string]string
//	@Failure		404	{object}	map[string]string
//	@Router			/private/tickets/{id} [get]
func GetTicketByID(c *gin.Context) {
	userID, err := services.ExtractTokenID(c)
	if err != nil {
		slog.Warn("Unauthorized ticket fetch attempt")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	ticketID := c.Param("id")

	ticket, err := services.GetTicketByID(ticketID, userID)
	if err != nil {
		slog.Warn("Ticket not found", "ticket_id", ticketID, "user_id", userID)
		c.JSON(http.StatusNotFound, gin.H{"error": "ticket not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": ticket})
}

// ValidateTicket godoc
//
//	@Summary		Validate a ticket
//	@Description	Mark a ticket as used — merchant only
//	@Tags			Tickets
//	@Produce		json
//	@Param			id	path		string	true	"Ticket UUID"
//	@Success		200	{object}	dto.TicketResponse
//	@Failure		400	{object}	map[string]string
//	@Failure		401	{object}	map[string]string
//	@Failure		404	{object}	map[string]string
//	@Router			/merchant/tickets/{id}/validate [post]
func ValidateTicket(c *gin.Context) {
	merchantID, err := services.ExtractTokenID(c)
	if err != nil {
		slog.Warn("Unauthorized ticket validation attempt")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	ticketID := c.Param("id")

	ticket, err := services.ValidateTicket(ticketID, merchantID)
	if err != nil {
		slog.Warn("Ticket validation failed",
			"ticket_id", ticketID,
			"merchant_id", merchantID,
			"error", err.Error(),
		)
		switch {
		case errors.Is(err, services.ErrTicketNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, services.ErrTicketAlreadyUsed):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to validate ticket"})
		}
		return
	}

	slog.Info("Ticket validated", "ticket_id", ticketID, "merchant_id", merchantID)
	c.JSON(http.StatusOK, gin.H{"data": ticket})
}
