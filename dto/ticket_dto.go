// dto/ticket_dto.go
package dto

import "time"

type TicketResponse struct {
	TicketID       string         `json:"ticket_id"`
	TransactionID  string         `json:"transaction_id"`
	UserID         string         `json:"user_id"`
	EventID        string         `json:"event_id"`
	TicketTypeID   string         `json:"ticket_type_id"`
	TicketTypeName string         `json:"ticket_type_name"`
	PricePaidCents int64          `json:"price_paid_cents"`
	Status         string         `json:"status"`
	IsGift         bool           `json:"is_gift"`
	GiftedBy       string         `json:"gifted_by,omitempty"`
	GiftedAt       *time.Time     `json:"gifted_at,omitempty"`
	GiftMessage    string         `json:"gift_message,omitempty"`
	Event          *EventResponse `json:"event,omitempty"`
	CreatedAt      time.Time      `json:"created_at"`
}

type PaginatedTicketsResponse struct {
	Data  []TicketResponse `json:"data"`
	Page  int              `json:"page"`
	Limit int              `json:"limit"`
	Total int64            `json:"total"`
}
