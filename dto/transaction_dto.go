// dto/transaction_dto.go
package dto

import "time"

type PurchaseTicketRequest struct {
	EventID string  `json:"event_id" validate:"required,uuid"`
	Amount  float64 `json:"amount"   validate:"required,gt=0"`
}

type RefundRequest struct {
	TransactionID string `json:"transaction_id" validate:"required,uuid"`
}

type TransactionResponse struct {
	TransactionID string    `json:"transaction_id"`
	UserID        string    `json:"user_id"`
	EventID       string    `json:"event_id"`
	Amount        float64   `json:"amount"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
}

type PaginatedTransactionsResponse struct {
	Data  []TransactionResponse `json:"data"`
	Page  int                   `json:"page"`
	Limit int                   `json:"limit"`
	Total int64                 `json:"total"`
}
