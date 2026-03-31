// dto/ticket_type_dto.go
package dto

import "time"

type CreateTicketTypeRequest struct {
	Name         string     `json:"name"          validate:"required,min=2,max=100"`
	Description  string     `json:"description"   validate:"omitempty,max=500"`
	Category     string     `json:"category"      validate:"required,oneof=general vip early_bird reserved group day_pass tiered complimentary demographic"`
	PriceCents   int64      `json:"price_cents"   validate:"min=0"` // 0 = free/complimentary
	Capacity     int        `json:"capacity"      validate:"required,min=1"`
	MinPerOrder  int        `json:"min_per_order" validate:"omitempty,min=1"`
	MaxPerOrder  int        `json:"max_per_order" validate:"omitempty,min=1"`
	SaleStartsAt *time.Time `json:"sale_starts_at" validate:"omitempty"`
	SaleEndsAt   *time.Time `json:"sale_ends_at"   validate:"omitempty"`
	SortOrder    int        `json:"sort_order"    validate:"omitempty"`
}

type UpdateTicketTypeRequest struct {
	Name         string     `json:"name"          validate:"omitempty,min=2,max=100"`
	Description  string     `json:"description"   validate:"omitempty,max=500"`
	PriceCents   *int64     `json:"price_cents"   validate:"omitempty,min=0"`
	Capacity     *int       `json:"capacity"      validate:"omitempty,min=1"`
	MinPerOrder  *int       `json:"min_per_order" validate:"omitempty,min=1"`
	MaxPerOrder  *int       `json:"max_per_order" validate:"omitempty,min=1"`
	SaleStartsAt *time.Time `json:"sale_starts_at" validate:"omitempty"`
	SaleEndsAt   *time.Time `json:"sale_ends_at"   validate:"omitempty"`
	Active       *bool      `json:"active"        validate:"omitempty"`
	SortOrder    *int       `json:"sort_order"    validate:"omitempty"`
}

type TicketTypeResponse struct {
	TicketTypeID   string     `json:"ticket_type_id"`
	EventID        string     `json:"event_id"`
	Name           string     `json:"name"`
	Description    string     `json:"description"`
	Category       string     `json:"category"`
	PriceCents     int64      `json:"price_cents"`
	PriceFormatted string     `json:"price_formatted"` // "R$ 50,00"
	Capacity       int        `json:"capacity"`
	Available      int        `json:"available"`
	Sold           int        `json:"sold"` // capacity - available
	MinPerOrder    int        `json:"min_per_order"`
	MaxPerOrder    int        `json:"max_per_order"`
	SaleStartsAt   *time.Time `json:"sale_starts_at"`
	SaleEndsAt     *time.Time `json:"sale_ends_at"`
	OnSale         bool       `json:"on_sale"` // computed: active + within sale window + available > 0
	Active         bool       `json:"active"`
	SortOrder      int        `json:"sort_order"`
	CreatedAt      time.Time  `json:"created_at"`
}
