// services/ticket_type_service.go
package services

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/maycolacerda/ticketfair/database"
	"github.com/maycolacerda/ticketfair/dto"
	"github.com/maycolacerda/ticketfair/models"
)

func CreateTicketType(merchantID, eventID string, req dto.CreateTicketTypeRequest) (*dto.TicketTypeResponse, error) {
	// Verify event belongs to merchant
	var event models.Event
	if err := database.DB.
		Where("event_id = ? AND merchant_id = ? AND deleted_at IS NULL", eventID, merchantID).
		First(&event).Error; err != nil {
		return nil, ErrEventNotFound
	}

	minPerOrder := req.MinPerOrder
	if minPerOrder == 0 {
		minPerOrder = 1
	}
	maxPerOrder := req.MaxPerOrder
	if maxPerOrder == 0 {
		maxPerOrder = 10
	}

	tt := models.TicketType{
		EventID:      eventID,
		Name:         strings.TrimSpace(req.Name),
		Description:  strings.TrimSpace(req.Description),
		Category:     req.Category,
		PriceCents:   req.PriceCents,
		Capacity:     req.Capacity,
		Available:    req.Capacity, // starts fully available
		MinPerOrder:  minPerOrder,
		MaxPerOrder:  maxPerOrder,
		SaleStartsAt: req.SaleStartsAt,
		SaleEndsAt:   req.SaleEndsAt,
		SortOrder:    req.SortOrder,
		Active:       true,
	}

	if err := database.DB.Create(&tt).Error; err != nil {
		return nil, ErrFailedToCreate
	}

	slog.Info("Ticket type created",
		"ticket_type_id", tt.TicketTypeID,
		"event_id", eventID,
		"category", req.Category,
	)

	return toTicketTypeResponse(&tt), nil
}

func UpdateTicketType(merchantID, eventID, ticketTypeID string, req dto.UpdateTicketTypeRequest) (*dto.TicketTypeResponse, error) {
	var tt models.TicketType

	// Scope to merchant via event
	if err := database.DB.
		Joins("JOIN events ON events.event_id = ticket_types.event_id").
		Where("ticket_types.ticket_type_id = ? AND ticket_types.event_id = ? AND events.merchant_id = ? AND ticket_types.deleted_at IS NULL",
			ticketTypeID, eventID, merchantID).
		First(&tt).Error; err != nil {
		return nil, ErrTicketTypeNotFound
	}

	updates := map[string]interface{}{}

	if req.Name != "" {
		updates["name"] = strings.TrimSpace(req.Name)
	}
	if req.Description != "" {
		updates["description"] = strings.TrimSpace(req.Description)
	}
	if req.PriceCents != nil {
		updates["price_cents"] = *req.PriceCents
	}
	if req.Capacity != nil {
		// Adjust available proportionally
		delta := *req.Capacity - tt.Capacity
		newAvailable := tt.Available + delta
		if newAvailable < 0 {
			newAvailable = 0
		}
		updates["capacity"] = *req.Capacity
		updates["available"] = newAvailable
	}
	if req.MinPerOrder != nil {
		updates["min_per_order"] = *req.MinPerOrder
	}
	if req.MaxPerOrder != nil {
		updates["max_per_order"] = *req.MaxPerOrder
	}
	if req.SaleStartsAt != nil {
		updates["sale_starts_at"] = req.SaleStartsAt
	}
	if req.SaleEndsAt != nil {
		updates["sale_ends_at"] = req.SaleEndsAt
	}
	if req.Active != nil {
		updates["active"] = *req.Active
	}
	if req.SortOrder != nil {
		updates["sort_order"] = *req.SortOrder
	}

	if len(updates) == 0 {
		return nil, ErrNoFieldsToUpdate
	}

	if err := database.DB.Model(&tt).Updates(updates).Error; err != nil {
		return nil, ErrFailedToUpdate
	}

	if err := database.DB.First(&tt, "ticket_type_id = ?", ticketTypeID).Error; err != nil {
		return nil, ErrFailedToFetch
	}

	slog.Info("Ticket type updated", "ticket_type_id", ticketTypeID)
	return toTicketTypeResponse(&tt), nil
}

func DeleteTicketType(merchantID, eventID, ticketTypeID string) error {
	var tt models.TicketType

	if err := database.DB.
		Joins("JOIN events ON events.event_id = ticket_types.event_id").
		Where("ticket_types.ticket_type_id = ? AND ticket_types.event_id = ? AND events.merchant_id = ? AND ticket_types.deleted_at IS NULL",
			ticketTypeID, eventID, merchantID).
		First(&tt).Error; err != nil {
		return ErrTicketTypeNotFound
	}

	// Soft delete only if no tickets sold
	if tt.Available < tt.Capacity {
		return ErrTicketTypeHasSales
	}

	if err := database.DB.Delete(&tt).Error; err != nil {
		return ErrFailedToUpdate
	}

	slog.Info("Ticket type deleted", "ticket_type_id", ticketTypeID)
	return nil
}

func GetTicketTypesByEvent(eventID string) ([]dto.TicketTypeResponse, error) {
	var types []models.TicketType

	if err := database.DB.
		Where("event_id = ? AND active = true AND deleted_at IS NULL", eventID).
		Order("sort_order ASC, created_at ASC").
		Find(&types).Error; err != nil {
		return nil, ErrFailedToFetch
	}

	result := make([]dto.TicketTypeResponse, len(types))
	for i, tt := range types {
		result[i] = *toTicketTypeResponse(&tt)
	}
	return result, nil
}

func GetAllTicketTypesByEvent(merchantID, eventID string) ([]dto.TicketTypeResponse, error) {
	// Merchant version — sees inactive types too
	var types []models.TicketType

	if err := database.DB.
		Joins("JOIN events ON events.event_id = ticket_types.event_id").
		Where("ticket_types.event_id = ? AND events.merchant_id = ? AND ticket_types.deleted_at IS NULL",
			eventID, merchantID).
		Order("ticket_types.sort_order ASC, ticket_types.created_at ASC").
		Find(&types).Error; err != nil {
		return nil, ErrFailedToFetch
	}

	result := make([]dto.TicketTypeResponse, len(types))
	for i, tt := range types {
		result[i] = *toTicketTypeResponse(&tt)
	}
	return result, nil
}

// ─────────────────────────────────────────────
// Updated PurchaseTicket — now ticket-type-aware
// ─────────────────────────────────────────────

func ComputeOrderAmount(ticketTypeID string, quantity int) (float64, error) {
	var tt models.TicketType
	if err := database.DB.First(&tt, "ticket_type_id = ?", ticketTypeID).Error; err != nil {
		return 0, ErrTicketTypeNotFound
	}
	return float64(tt.PriceCents*int64(quantity)) / 100, nil
}

// ─────────────────────────────────────────────
// Helpers
// ─────────────────────────────────────────────

func toTicketTypeResponse(tt *models.TicketType) *dto.TicketTypeResponse {
	now := time.Now()

	onSale := tt.Active &&
		tt.Available > 0 &&
		(tt.SaleStartsAt == nil || now.After(*tt.SaleStartsAt)) &&
		(tt.SaleEndsAt == nil || now.Before(*tt.SaleEndsAt))

	return &dto.TicketTypeResponse{
		TicketTypeID:   tt.TicketTypeID,
		EventID:        tt.EventID,
		Name:           tt.Name,
		Description:    tt.Description,
		Category:       tt.Category,
		PriceCents:     tt.PriceCents,
		PriceFormatted: formatBRL(tt.PriceCents),
		Capacity:       tt.Capacity,
		Available:      tt.Available,
		Sold:           tt.Capacity - tt.Available,
		MinPerOrder:    tt.MinPerOrder,
		MaxPerOrder:    tt.MaxPerOrder,
		SaleStartsAt:   tt.SaleStartsAt,
		SaleEndsAt:     tt.SaleEndsAt,
		OnSale:         onSale,
		Active:         tt.Active,
		SortOrder:      tt.SortOrder,
		CreatedAt:      tt.CreatedAt,
	}
}

func formatBRL(cents int64) string {
	if cents == 0 {
		return "Free"
	}
	reais := cents / 100
	centavos := cents % 100
	return fmt.Sprintf("R$ %d,%02d", reais, centavos)
}
