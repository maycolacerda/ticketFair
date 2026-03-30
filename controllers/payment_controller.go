// controllers/payment.go
package controllers

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/maycolacerda/ticketfair/dto"
	"github.com/maycolacerda/ticketfair/services"
)

// CreatePaymentIntent godoc
//
//	@Summary		Create a Stripe PaymentIntent
//	@Description	Returns a client_secret to confirm payment on the client side
//	@Tags			Payments
//	@Accept			json
//	@Produce		json
//	@Param			body	body		dto.CreatePaymentIntentRequest	true	"Event and quantity"
//	@Success		201		{object}	dto.CreatePaymentIntentResponse
//	@Failure		400		{object}	map[string]string
//	@Failure		401		{object}	map[string]string
//	@Failure		404		{object}	map[string]string
//	@Failure		409		{object}	map[string]string
//	@Failure		503		{object}	map[string]string
//	@Router			/private/payments/intent [post]
func CreatePaymentIntent(c *gin.Context) {
	userID, err := services.ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req dto.CreatePaymentIntentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if err := validate.Struct(req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": formatValidationErrors(err)})
		return
	}

	resp, err := services.CreatePaymentIntent(userID, req)
	if err != nil {
		slog.Warn("CreatePaymentIntent failed",
			"user_id", userID,
			"error", err.Error(),
		)
		switch {
		case errors.Is(err, services.ErrStripeDisabled):
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		case errors.Is(err, services.ErrEventNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, services.ErrEventSoldOut):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		case errors.Is(err, services.ErrUserNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create payment"})
		}
		return
	}

	slog.Info("PaymentIntent created",
		"user_id", userID,
		"event_id", req.EventID,
		"pi_id", resp.PaymentIntentID,
	)
	c.JSON(http.StatusCreated, gin.H{"data": resp})
}

// StripeWebhook godoc
//
//	@Summary		Stripe webhook receiver
//	@Description	Receives and processes Stripe webhook events (payment_intent.succeeded, charge.refunded, etc.)
//	@Tags			Payments
//	@Accept			plain
//	@Produce		json
//	@Success		200	{object}	dto.WebhookEventResponse
//	@Failure		400	{object}	map[string]string
//	@Router			/public/webhooks/stripe [post]
func StripeWebhook(c *gin.Context) {
	// Read raw body — required for signature verification
	payload, err := io.ReadAll(c.Request.Body)
	if err != nil {
		slog.Warn("Webhook: failed to read body", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read request body"})
		return
	}

	sigHeader := c.GetHeader("Stripe-Signature")
	if sigHeader == "" {
		slog.Warn("Webhook: missing Stripe-Signature header")
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing Stripe-Signature header"})
		return
	}

	resp, err := services.HandleWebhook(payload, sigHeader)
	if err != nil {
		slog.Warn("Webhook processing failed", "error", err.Error())
		switch {
		case errors.Is(err, services.ErrInvalidWebhook):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid webhook signature"})
		default:
			// Return 200 anyway — Stripe will retry on 4xx/5xx
			// but we don't want infinite retries for internal errors
			c.JSON(http.StatusOK, gin.H{"received": true})
		}
		return
	}

	c.JSON(http.StatusOK, resp)
}

// RefundPayment godoc
//
//	@Summary		Refund a payment
//	@Description	Issues a full Stripe refund for a succeeded payment
//	@Tags			Payments
//	@Produce		json
//	@Param			id	path		string	true	"Payment UUID"
//	@Success		200	{object}	map[string]string
//	@Failure		400	{object}	map[string]string
//	@Failure		401	{object}	map[string]string
//	@Failure		404	{object}	map[string]string
//	@Router			/private/payments/{id}/refund [post]
func RefundPayment(c *gin.Context) {
	userID, err := services.ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	paymentID := c.Param("id")

	if err := services.RefundPayment(userID, paymentID); err != nil {
		slog.Warn("Refund failed",
			"payment_id", paymentID,
			"user_id", userID,
			"error", err.Error(),
		)
		switch {
		case errors.Is(err, services.ErrStripeDisabled):
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		case errors.Is(err, services.ErrPaymentNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, services.ErrPaymentNotSucceeded),
			errors.Is(err, services.ErrNotRefundable):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to issue refund"})
		}
		return
	}

	slog.Info("Refund issued", "payment_id", paymentID, "user_id", userID)
	c.JSON(http.StatusOK, gin.H{"message": "refund issued successfully"})
}

// GetMyPayments godoc
//
//	@Summary		List user payments
//	@Description	Paginated payment history for the authenticated user
//	@Tags			Payments
//	@Produce		json
//	@Param			page	query		int	false	"Page"	default(1)
//	@Param			limit	query		int	false	"Limit"	default(20)
//	@Success		200		{object}	dto.PaginatedPaymentsResponse
//	@Failure		401		{object}	map[string]string
//	@Router			/private/payments [get]
func GetMyPayments(c *gin.Context) {
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

	result, err := services.GetUserPayments(userID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch payments"})
		return
	}

	c.JSON(http.StatusOK, result)
}
