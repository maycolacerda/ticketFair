// controllers/transaction.go
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

// PurchaseTicket godoc
//
//	@Summary		Purchase a ticket
//	@Description	Purchase a ticket for an event
//	@Tags			Transactions
//	@Accept			json
//	@Produce		json
//	@Param			purchase	body		dto.PurchaseTicketRequest	true	"Purchase data"
//	@Success		201			{object}	dto.TransactionResponse
//	@Failure		400			{object}	map[string]string
//	@Failure		401			{object}	map[string]string
//	@Failure		404			{object}	map[string]string
//	@Failure		409			{object}	map[string]string
//	@Failure		422			{object}	map[string]interface{}
//	@Router			/private/tickets/purchase [post]
func PurchaseTicket(c *gin.Context) {
	userID, err := services.ExtractTokenID(c)
	if err != nil {
		slog.Warn("Unauthorized purchase attempt")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req dto.PurchaseTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Warn("Invalid request body", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := validate.Struct(req); err != nil {
		errs := formatValidationErrors(err)
		slog.Warn("Purchase validation failed", "errors", errs)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": errs})
		return
	}

	tx, err := services.PurchaseTicket(userID, req.EventID, req.Amount)
	if err != nil {
		slog.Warn("Purchase failed",
			"user_id", userID,
			"event_id", req.EventID,
			"error", err.Error(),
		)
		switch {
		case errors.Is(err, services.ErrEventNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, services.ErrEventSoldOut):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to purchase ticket"})
		}
		return
	}

	slog.Info("Purchase successful",
		"transaction_id", tx.TransactionID,
		"user_id", userID,
	)
	c.JSON(http.StatusCreated, gin.H{"data": tx})
}

// RefundTicket godoc
//
//	@Summary		Refund a ticket
//	@Description	Refund a completed ticket transaction
//	@Tags			Transactions
//	@Accept			json
//	@Produce		json
//	@Param			refund	body		dto.RefundRequest	true	"Refund data"
//	@Success		200		{object}	map[string]string
//	@Failure		400		{object}	map[string]string
//	@Failure		401		{object}	map[string]string
//	@Failure		404		{object}	map[string]string
//	@Failure		422		{object}	map[string]interface{}
//	@Router			/private/tickets/refund [post]
func RefundTicket(c *gin.Context) {
	userID, err := services.ExtractTokenID(c)
	if err != nil {
		slog.Warn("Unauthorized refund attempt")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req dto.RefundRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Warn("Invalid request body", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := validate.Struct(req); err != nil {
		errs := formatValidationErrors(err)
		slog.Warn("Refund validation failed", "errors", errs)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": errs})
		return
	}

	if err := services.RefundTicket(req.TransactionID, userID); err != nil {
		slog.Warn("Refund failed",
			"transaction_id", req.TransactionID,
			"user_id", userID,
			"error", err.Error(),
		)
		switch {
		case errors.Is(err, services.ErrTransactionNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, services.ErrNotRefundable):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process refund"})
		}
		return
	}

	slog.Info("Refund successful",
		"transaction_id", req.TransactionID,
		"user_id", userID,
	)
	c.JSON(http.StatusOK, gin.H{"message": "ticket refunded successfully"})
}

// GetMyTransactions godoc
//
//	@Summary		List user transactions
//	@Description	Retrieve paginated transaction history for the authenticated user
//	@Tags			Transactions
//	@Produce		json
//	@Param			page	query		int	false	"Page number"	default(1)
//	@Param			limit	query		int	false	"Page size"		default(20)
//	@Success		200		{object}	dto.PaginatedTransactionsResponse
//	@Failure		401		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/private/transactions [get]
func GetMyTransactions(c *gin.Context) {
	userID, err := services.ExtractTokenID(c)
	if err != nil {
		slog.Warn("Unauthorized transaction list attempt")
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

	result, err := services.GetUserTransactions(userID, page, limit)
	if err != nil {
		slog.Error("Failed to fetch transactions",
			"user_id", userID,
			"error", err.Error(),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch transactions"})
		return
	}

	c.JSON(http.StatusOK, result)
}
