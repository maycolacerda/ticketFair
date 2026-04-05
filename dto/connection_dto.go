// dto/connection_dto.go
package dto

import "time"

// ── Connections ───────────────────────────────────────

type SendConnectionRequest struct {
	AddresseeID string `json:"addressee_id" validate:"required,uuid"`
}

type RespondConnectionRequest struct {
	Action string `json:"action" validate:"required,oneof=accept decline"`
}

type ConnectionResponse struct {
	ConnectionID string      `json:"connection_id"`
	Status       string      `json:"status"`
	User         UserSummary `json:"user"` // the other person
	CreatedAt    time.Time   `json:"created_at"`
}

type UserSummary struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type PaginatedConnectionsResponse struct {
	Data  []ConnectionResponse `json:"data"`
	Page  int                  `json:"page"`
	Limit int                  `json:"limit"`
	Total int64                `json:"total"`
}

// ── Gift Tickets ──────────────────────────────────────

type GiftTicketRequest struct {
	RecipientID string `json:"recipient_id" validate:"required,uuid"`
	Message     string `json:"message"      validate:"omitempty,max=300"`
}

type GiftTicketResponse struct {
	TicketID    string    `json:"ticket_id"`
	EventName   string    `json:"event_name"`
	RecipientID string    `json:"recipient_id"`
	Recipient   string    `json:"recipient_username"`
	Message     string    `json:"message"`
	GiftedAt    time.Time `json:"gifted_at"`
}

// ── Connection Events feed ────────────────────────────

type ConnectionEventResponse struct {
	EventID        string        `json:"event_id"`
	Name           string        `json:"name"`
	Location       string        `json:"location"`
	StartTime      time.Time     `json:"start_time"`
	ImageURL       string        `json:"image_url"`
	Attendees      []UserSummary `json:"attendees"` // connections attending
	AttendeesCount int           `json:"attendees_count"`
}
