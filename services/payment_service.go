// services/payment_service.go
package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"

	configs "github.com/maycolacerda/ticketfair/configs/stripe"
	"github.com/maycolacerda/ticketfair/database"
	"github.com/maycolacerda/ticketfair/dto"
	"github.com/maycolacerda/ticketfair/models"
	"github.com/stripe/stripe-go/v79"
	"github.com/stripe/stripe-go/v79/paymentintent"
	"github.com/stripe/stripe-go/v79/refund"
	"github.com/stripe/stripe-go/v79/webhook"
)

// ── Errors ────────────────────────────────────────────

var (
	ErrStripeDisabled      = errors.New("payment processing is not configured")
	ErrPaymentNotFound     = errors.New("payment not found")
	ErrPaymentNotSucceeded = errors.New("payment has not succeeded")
	ErrPaymentAlreadyDone  = errors.New("payment already processed")
	ErrInvalidWebhook      = errors.New("invalid webhook signature")
)

const defaultTicketPriceCents = 5000 // R$ 50,00

func CreatePaymentIntent(userID string, req dto.CreatePaymentIntentRequest) (*dto.CreatePaymentIntentResponse, error) {
	if !configs.StripeEnabled() {
		return nil, ErrStripeDisabled
	}

	// Verify event exists and has capacity
	var event models.Event
	if err := database.DB.First(&event, "event_id = ? AND active = true AND deleted_at IS NULL", req.EventID).Error; err != nil {
		return nil, ErrEventNotFound
	}
	if event.Capacity < req.Quantity {
		return nil, ErrEventSoldOut
	}

	// Verify user exists
	var user models.User
	if err := database.DB.First(&user, "user_id = ? AND active = true", userID).Error; err != nil {
		return nil, ErrUserNotFound
	}

	totalCents := int64(defaultTicketPriceCents) * int64(req.Quantity)

	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(totalCents),
		Currency: stripe.String("brl"),
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
		Metadata: map[string]string{
			"user_id":  userID,
			"event_id": req.EventID,
			"quantity": fmt.Sprintf("%d", req.Quantity),
		},
		Description:  stripe.String(fmt.Sprintf("TicketFair — %s x%d", event.Name, req.Quantity)),
		ReceiptEmail: stripe.String(user.Email),
	}

	pi, err := paymentintent.New(params)
	if err != nil {
		slog.Error("Stripe PaymentIntent creation failed",
			"user_id", userID,
			"event_id", req.EventID,
			"error", err.Error(),
		)
		return nil, ErrFailedToCreate
	}

	// Persist pending payment record
	payment := models.Payment{
		UserID:          userID,
		EventID:         req.EventID,
		StripePaymentID: pi.ID,
		Amount:          totalCents,
		Currency:        "brl",
		Quantity:        req.Quantity,
		Status:          "pending",
	}
	if err := database.DB.Create(&payment).Error; err != nil {

		_, _ = paymentintent.Cancel(pi.ID, nil)
		return nil, ErrFailedToCreate
	}

	slog.Info("PaymentIntent created",
		"payment_intent_id", pi.ID,
		"user_id", userID,
		"event_id", req.EventID,
		"amount_cents", totalCents,
	)

	return &dto.CreatePaymentIntentResponse{
		ClientSecret:    pi.ClientSecret,
		PaymentIntentID: pi.ID,
		Amount:          totalCents,
		Currency:        "brl",
		EventID:         req.EventID,
		EventName:       event.Name,
		Quantity:        req.Quantity,
	}, nil
}

func HandleWebhook(payload []byte, sigHeader string) (*dto.WebhookEventResponse, error) {
	webhookSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")
	if webhookSecret == "" {
		return nil, ErrInvalidWebhook
	}

	event, err := webhook.ConstructEvent(payload, sigHeader, webhookSecret)
	if err != nil {
		slog.Warn("Webhook signature verification failed", "error", err.Error())
		return nil, ErrInvalidWebhook
	}

	slog.Info("Webhook received", "type", event.Type, "id", event.ID)

	switch event.Type {

	case "payment_intent.succeeded":
		return handlePaymentSucceeded(event)

	case "payment_intent.payment_failed":
		return handlePaymentFailed(event)

	case "payment_intent.canceled":
		return handlePaymentCanceled(event)

	case "charge.refunded":
		return handleChargeRefunded(event)

	default:
		slog.Info("Unhandled webhook event", "type", event.Type)
		return &dto.WebhookEventResponse{Received: true}, nil
	}
}

func handlePaymentSucceeded(event stripe.Event) (*dto.WebhookEventResponse, error) {
	var pi stripe.PaymentIntent
	if err := json.Unmarshal(event.Data.Raw, &pi); err != nil {
		return nil, ErrFailedToFetch
	}

	// Find pending payment record
	var payment models.Payment
	if err := database.DB.Where("stripe_payment_id = ?", pi.ID).First(&payment).Error; err != nil {
		slog.Error("Payment record not found for succeeded intent", "pi_id", pi.ID)
		return nil, ErrPaymentNotFound
	}

	// Idempotency check
	if payment.Status == "succeeded" {
		slog.Warn("PaymentIntent already processed", "pi_id", pi.ID)
		return &dto.WebhookEventResponse{Received: true, EventID: event.ID}, nil
	}

	var transactionID string
	err := database.DB.Raw(
		`SELECT purchase_ticket(?, ?, ?)`,
		payment.UserID, payment.EventID, float64(payment.Amount)/100,
	).Scan(&transactionID).Error

	if err != nil {
		slog.Error("purchase_ticket failed after payment success",
			"pi_id", pi.ID,
			"user_id", payment.UserID,
			"error", err.Error(),
		)
		// Mark payment as failed to allow retry handling
		database.DB.Model(&payment).Update("status", "failed")
		return nil, ErrFailedToCreate
	}

	_, ticketErr := CreateTicket(transactionID, payment.UserID, payment.EventID)
	if ticketErr != nil {
		slog.Error("Failed to create ticket after payment",
			"transaction_id", transactionID,
			"error", ticketErr.Error(),
		)
	}

	database.DB.Model(&payment).Updates(map[string]interface{}{
		"status":         "succeeded",
		"transaction_id": transactionID,
	})

	go func() {
		var user models.User
		var ev models.Event
		if err := database.DB.First(&user, "user_id = ?", payment.UserID).Error; err != nil {
			return
		}
		if err := database.DB.First(&ev, "event_id = ?", payment.EventID).Error; err != nil {
			return
		}
		_ = SendPurchaseConfirmationEmail(
			user.Email,
			user.Username,
			ev.Name,
			transactionID,
			float64(payment.Amount)/100,
		)
	}()

	slog.Info("Payment succeeded — ticket created",
		"pi_id", pi.ID,
		"transaction_id", transactionID,
		"user_id", payment.UserID,
		"event_id", payment.EventID,
	)

	return &dto.WebhookEventResponse{Received: true, EventID: event.ID}, nil
}

func handlePaymentFailed(event stripe.Event) (*dto.WebhookEventResponse, error) {
	var pi stripe.PaymentIntent
	if err := json.Unmarshal(event.Data.Raw, &pi); err != nil {
		return nil, ErrFailedToFetch
	}

	database.DB.Model(&models.Payment{}).
		Where("stripe_payment_id = ?", pi.ID).
		Update("status", "failed")

	slog.Warn("Payment failed", "pi_id", pi.ID)
	return &dto.WebhookEventResponse{Received: true, EventID: event.ID}, nil
}

func handlePaymentCanceled(event stripe.Event) (*dto.WebhookEventResponse, error) {
	var pi stripe.PaymentIntent
	if err := json.Unmarshal(event.Data.Raw, &pi); err != nil {
		return nil, ErrFailedToFetch
	}

	database.DB.Model(&models.Payment{}).
		Where("stripe_payment_id = ?", pi.ID).
		Update("status", "canceled")

	slog.Info("Payment canceled", "pi_id", pi.ID)
	return &dto.WebhookEventResponse{Received: true, EventID: event.ID}, nil
}

func handleChargeRefunded(event stripe.Event) (*dto.WebhookEventResponse, error) {
	var charge stripe.Charge
	if err := json.Unmarshal(event.Data.Raw, &charge); err != nil {
		return nil, ErrFailedToFetch
	}

	piID := ""
	if charge.PaymentIntent != nil {
		piID = charge.PaymentIntent.ID
	}
	if piID == "" {
		return &dto.WebhookEventResponse{Received: true}, nil
	}

	var payment models.Payment
	if err := database.DB.Where("stripe_payment_id = ?", piID).First(&payment).Error; err != nil {
		return &dto.WebhookEventResponse{Received: true}, nil
	}

	// Run the existing refund SQL function to restore capacity
	if payment.TransactionID != "" {
		_ = database.DB.Exec(`SELECT refund_ticket(?)`, payment.TransactionID).Error
		database.DB.Model(&models.Ticket{}).
			Where("transaction_id = ?", payment.TransactionID).
			Update("status", "refunded")
	}

	database.DB.Model(&payment).Update("status", "refunded")

	slog.Info("Charge refunded",
		"pi_id", piID,
		"transaction_id", payment.TransactionID,
	)

	return &dto.WebhookEventResponse{Received: true, EventID: event.ID}, nil
}

func RefundPayment(userID, paymentID string) error {
	if !configs.StripeEnabled() {
		return ErrStripeDisabled
	}

	var payment models.Payment
	if err := database.DB.
		Where("payment_id = ? AND user_id = ?", paymentID, userID).
		First(&payment).Error; err != nil {
		return ErrPaymentNotFound
	}

	if payment.Status != "succeeded" {
		return ErrPaymentNotSucceeded
	}

	// Check ticket hasn't been used
	if payment.TransactionID != "" {
		var ticket models.Ticket
		if err := database.DB.Where("transaction_id = ?", payment.TransactionID).First(&ticket).Error; err == nil {
			if ticket.Status == "used" {
				return ErrNotRefundable
			}
		}
	}

	params := &stripe.RefundParams{
		PaymentIntent: stripe.String(payment.StripePaymentID),
	}
	_, err := refund.New(params)
	if err != nil {
		slog.Error("Stripe refund failed",
			"payment_id", paymentID,
			"pi_id", payment.StripePaymentID,
			"error", err.Error(),
		)
		return ErrFailedToUpdate
	}

	database.DB.Model(&payment).Update("status", "refunded")

	slog.Info("Refund issued",
		"payment_id", paymentID,
		"pi_id", payment.StripePaymentID,
		"user_id", userID,
	)
	return nil
}

func GetUserPayments(userID string, page, limit int) (*dto.PaginatedPaymentsResponse, error) {
	var payments []models.Payment
	var total int64
	offset := (page - 1) * limit

	if err := database.DB.Model(&models.Payment{}).
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Count(&total).Error; err != nil {
		return nil, ErrFailedToFetch
	}

	if err := database.DB.
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&payments).Error; err != nil {
		return nil, ErrFailedToFetch
	}

	data := make([]dto.PaymentResponse, len(payments))
	for i, p := range payments {
		data[i] = toPaymentResponse(&p)
	}

	return &dto.PaginatedPaymentsResponse{
		Data:  data,
		Page:  page,
		Limit: limit,
		Total: total,
	}, nil
}

func toPaymentResponse(p *models.Payment) dto.PaymentResponse {
	return dto.PaymentResponse{
		PaymentID:       p.PaymentID,
		UserID:          p.UserID,
		EventID:         p.EventID,
		TransactionID:   p.TransactionID,
		StripePaymentID: p.StripePaymentID,
		Amount:          p.Amount,
		Currency:        p.Currency,
		Status:          p.Status,
		CreatedAt:       p.CreatedAt,
	}
}
