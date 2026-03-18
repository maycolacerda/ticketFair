// dto/event_dto.go
package dto

import "time"

type CreateEventRequest struct {
	Name        string    `json:"name"        validate:"required,min=2,max=100"`
	Description string    `json:"description" validate:"max=1000"`
	Location    string    `json:"location"    validate:"required"`
	StartTime   time.Time `json:"start_time"  validate:"required"`
	EndTime     time.Time `json:"end_time"    validate:"required,gtfield=StartTime"`
	Capacity    int       `json:"capacity"    validate:"required,min=1"`
}

type UpdateEventRequest struct {
	Name        string    `json:"name"        validate:"omitempty,min=2,max=100"`
	Description string    `json:"description" validate:"omitempty,max=1000"`
	Location    string    `json:"location"    validate:"omitempty"`
	StartTime   time.Time `json:"start_time"  validate:"omitempty"`
	EndTime     time.Time `json:"end_time"    validate:"omitempty,gtfield=StartTime"`
	Capacity    int       `json:"capacity"    validate:"omitempty,min=1"`
	Active      *bool     `json:"active"      validate:"omitempty"` // pointer — false is a valid value
}

type EventResponse struct {
	EventID     string    `json:"event_id"`
	MerchantID  string    `json:"merchant_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Capacity    int       `json:"capacity"`
	Active      bool      `json:"active"`
	CreatedAt   time.Time `json:"created_at"`
}

type PaginatedEventsResponse struct {
	Data  []EventResponse `json:"data"`
	Page  int             `json:"page"`
	Limit int             `json:"limit"`
	Total int64           `json:"total"`
}
