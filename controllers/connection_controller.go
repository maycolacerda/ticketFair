// controllers/connection.go
package controllers

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/maycolacerda/ticketfair/dto"
	"github.com/maycolacerda/ticketfair/services"
)

// SendConnectionRequest godoc
//
//	@Summary		Send a connection request
//	@Description	Send a friend/connection request to another user
//	@Tags			Connections
//	@Accept			json
//	@Produce		json
//	@Param			body	body		dto.SendConnectionRequest	true	"Addressee"
//	@Success		201		{object}	dto.ConnectionResponse
//	@Router			/private/connections [post]
func SendConnectionRequest(c *gin.Context) {
	requesterID, err := services.ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req dto.SendConnectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if err := validate.Struct(req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": formatValidationErrors(err)})
		return
	}

	conn, err := services.SendConnectionRequest(requesterID, req.AddresseeID)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrConnectionSelf):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, services.ErrUserNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, services.ErrConnectionExists):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to send request"})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": conn})
}

// RespondToConnection godoc
//
//	@Summary		Accept or decline a connection request
//	@Tags			Connections
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string							true	"Connection UUID"
//	@Param			body	body		dto.RespondConnectionRequest	true	"Action"
//	@Success		200		{object}	dto.ConnectionResponse
//	@Router			/private/connections/{id}/respond [post]
func RespondToConnection(c *gin.Context) {
	addresseeID, err := services.ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	connectionID := c.Param("id")

	var req dto.RespondConnectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if err := validate.Struct(req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": formatValidationErrors(err)})
		return
	}

	conn, err := services.RespondToConnection(addresseeID, connectionID, req.Action)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrConnectionNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, services.ErrConnectionNotPending):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to respond"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": conn})
}

// RemoveConnection godoc
//
//	@Summary		Remove a connection
//	@Tags			Connections
//	@Produce		json
//	@Param			id	path		string	true	"Connection UUID"
//	@Success		200	{object}	map[string]string
//	@Router			/private/connections/{id} [delete]
func RemoveConnection(c *gin.Context) {
	userID, err := services.ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	connectionID := c.Param("id")

	if err := services.RemoveConnection(userID, connectionID); err != nil {
		if errors.Is(err, services.ErrConnectionNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to remove connection"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "connection removed"})
}

// ListConnections godoc
//
//	@Summary		List connections
//	@Description	Returns accepted connections. Use ?status=pending for incoming requests.
//	@Tags			Connections
//	@Produce		json
//	@Param			status	query		string	false	"Filter by status (accepted|pending|declined)"
//	@Param			page	query		int		false	"Page"	default(1)
//	@Param			limit	query		int		false	"Limit"	default(20)
//	@Success		200		{object}	dto.PaginatedConnectionsResponse
//	@Router			/private/connections [get]
func ListConnections(c *gin.Context) {
	userID, err := services.ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	status := c.Query("status")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	result, err := services.GetConnections(userID, status, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch connections"})
		return
	}

	c.JSON(http.StatusOK, result)
}

// ListPendingRequests godoc
//
//	@Summary		List incoming pending connection requests
//	@Tags			Connections
//	@Produce		json
//	@Success		200	{array}		dto.ConnectionResponse
//	@Router			/private/connections/requests [get]
func ListPendingRequests(c *gin.Context) {
	userID, err := services.ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	data, err := services.GetPendingRequests(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch requests"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": data})
}

// GetConnectionEvents godoc
//
//	@Summary		Events my connections are attending
//	@Description	Returns upcoming events where at least one connection has an active ticket
//	@Tags			Connections
//	@Produce		json
//	@Param			page	query		int	false	"Page"	default(1)
//	@Param			limit	query		int	false	"Limit"	default(20)
//	@Success		200		{array}		dto.ConnectionEventResponse
//	@Router			/private/connections/events [get]
func GetConnectionEvents(c *gin.Context) {
	userID, err := services.ExtractTokenID(c)
	if err != nil {
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

	result, err := services.GetConnectionEvents(userID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch events"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}

// GiftTicket godoc
//
//	@Summary		Gift a ticket to a connection
//	@Description	Transfer an active ticket to a connected user. Only accepted connections can receive gifts.
//	@Tags			Connections
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string				true	"Ticket UUID"
//	@Param			body	body		dto.GiftTicketRequest	true	"Recipient and optional message"
//	@Success		200		{object}	dto.GiftTicketResponse
//	@Router			/private/tickets/{id}/gift [post]
func GiftTicket(c *gin.Context) {
	senderID, err := services.ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	ticketID := c.Param("id")

	var req dto.GiftTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if err := validate.Struct(req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": formatValidationErrors(err)})
		return
	}

	result, err := services.GiftTicket(senderID, ticketID, req)
	if err != nil {
		slog.Warn("GiftTicket failed",
			"sender_id", senderID,
			"ticket_id", ticketID,
			"error", err.Error(),
		)
		switch {
		case errors.Is(err, services.ErrTicketNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, services.ErrUserNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "recipient not found"})
		case errors.Is(err, services.ErrCannotGiftOwn):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, services.ErrNotConnected):
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		case errors.Is(err, services.ErrTicketNotActive):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, services.ErrTicketAlreadyGifted):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to gift ticket"})
		}
		return
	}

	slog.Info("Ticket gifted",
		"ticket_id", ticketID,
		"sender_id", senderID,
		"recipient_id", req.RecipientID,
	)
	c.JSON(http.StatusOK, gin.H{"data": result})
}
