// dto/transaction_dto.go
package dto

import "time"

type RefundRequest struct {
	TransactionID string `json:"transaction_id" validate:"required,uuid"`
}

type TransactionResponse struct {
	TransactionID  string    `json:"transaction_id"`
	UserID         string    `json:"user_id"`
	EventID        string    `json:"event_id"`
	TicketTypeID   string    `json:"ticket_type_id"`
	TicketTypeName string    `json:"ticket_type_name"`
	Quantity       int       `json:"quantity"`
	Amount         float64   `json:"amount"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
}

type PaginatedTransactionsResponse struct {
	Data  []TransactionResponse `json:"data"`
	Page  int                   `json:"page"`
	Limit int                   `json:"limit"`
	Total int64                 `json:"total"`
}

type PurchaseTicketRequest struct {
	EventID      string `json:"event_id"       validate:"required,uuid"`
	TicketTypeID string `json:"ticket_type_id" validate:"required,uuid"`
	Quantity     int    `json:"quantity"       validate:"required,min=1,max=10"`
}

type CreatePaymentIntentRequest struct {
	EventID      string `json:"event_id"       validate:"required,uuid"`
	TicketTypeID string `json:"ticket_type_id" validate:"required,uuid"`
	Quantity     int    `json:"quantity"       validate:"required,min=1,max=10"`
}
