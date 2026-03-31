// dto/payment_dto.go
package dto

import "time"

// ── Create Payment Intent ──────────────────────────────

type CreatePaymentIntentResponse struct {
	ClientSecret    string `json:"client_secret"`
	PaymentIntentID string `json:"payment_intent_id"`
	Amount          int64  `json:"amount"` // in cents
	Currency        string `json:"currency"`
	EventID         string `json:"event_id"`
	EventName       string `json:"event_name"`
	Quantity        int    `json:"quantity"`
}

// ── Confirm Payment ────────────────────────────────────

type ConfirmPaymentRequest struct {
	PaymentIntentID string `json:"payment_intent_id" validate:"required"`
}

// ── Payment record ─────────────────────────────────────

type PaymentResponse struct {
	PaymentID       string    `json:"payment_id"`
	UserID          string    `json:"user_id"`
	EventID         string    `json:"event_id"`
	TransactionID   string    `json:"transaction_id,omitempty"`
	StripePaymentID string    `json:"stripe_payment_id"`
	Amount          int64     `json:"amount"` // in cents
	Currency        string    `json:"currency"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
}

type PaginatedPaymentsResponse struct {
	Data  []PaymentResponse `json:"data"`
	Page  int               `json:"page"`
	Limit int               `json:"limit"`
	Total int64             `json:"total"`
}

// ── Webhook ────────────────────────────────────────────

type WebhookEventResponse struct {
	Received bool   `json:"received"`
	EventID  string `json:"event_id,omitempty"`
}
