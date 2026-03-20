// dto/ticket_dto.go
package dto

import "time"

type TicketResponse struct {
	TicketID      string         `json:"ticket_id"`
	TransactionID string         `json:"transaction_id"`
	UserID        string         `json:"user_id"`
	EventID       string         `json:"event_id"`
	Status        string         `json:"status"`
	Event         *EventResponse `json:"event,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
}

type PaginatedTicketsResponse struct {
	Data  []TicketResponse `json:"data"`
	Page  int              `json:"page"`
	Limit int              `json:"limit"`
	Total int64            `json:"total"`
}
