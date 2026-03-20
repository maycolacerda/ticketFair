// services/ticket_service.go
package services

import (
	"errors"
	"log/slog"

	"github.com/maycolacerda/ticketfair/database"
	"github.com/maycolacerda/ticketfair/dto"
	"github.com/maycolacerda/ticketfair/models"
)

var ErrTicketNotFound = errors.New("ticket not found")
var ErrTicketAlreadyUsed = errors.New("ticket already used")

func CreateTicket(transactionID, userID, eventID string) (*dto.TicketResponse, error) {
	ticket := models.Ticket{
		TransactionID: transactionID,
		UserID:        userID,
		EventID:       eventID,
		Status:        "active",
	}

	if err := database.DB.Create(&ticket).Error; err != nil {
		slog.Error("Failed to create ticket",
			"transaction_id", transactionID,
			"user_id", userID,
			"error", err.Error(),
		)
		return nil, ErrFailedToCreate
	}

	slog.Info("Ticket created",
		"ticket_id", ticket.TicketID,
		"user_id", userID,
		"event_id", eventID,
	)

	return toTicketResponse(&ticket, nil), nil
}

func GetTicketByID(ticketID, userID string) (*dto.TicketResponse, error) {
	var ticket models.Ticket

	if err := database.DB.
		Preload("Event").
		Where("ticket_id = ? AND user_id = ?", ticketID, userID).
		First(&ticket).Error; err != nil {
		return nil, ErrTicketNotFound
	}

	var eventResp *dto.EventResponse
	if ticket.Event != nil {
		eventResp = toEventResponse(ticket.Event)
	}

	return toTicketResponse(&ticket, eventResp), nil
}

func GetUserTickets(userID string, page, limit int) (*dto.PaginatedTicketsResponse, error) {
	var tickets []models.Ticket
	var total int64

	offset := (page - 1) * limit

	if err := database.DB.Model(&models.Ticket{}).
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Count(&total).Error; err != nil {
		return nil, ErrFailedToFetch
	}

	if err := database.DB.
		Preload("Event").
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&tickets).Error; err != nil {
		return nil, ErrFailedToFetch
	}

	data := make([]dto.TicketResponse, len(tickets))
	for i, t := range tickets {
		var eventResp *dto.EventResponse
		if t.Event != nil {
			eventResp = toEventResponse(t.Event)
		}
		data[i] = *toTicketResponse(&t, eventResp)
	}

	return &dto.PaginatedTicketsResponse{
		Data:  data,
		Page:  page,
		Limit: limit,
		Total: total,
	}, nil
}

func ValidateTicket(ticketID, merchantID string) (*dto.TicketResponse, error) {
	var ticket models.Ticket

	// Preload event to verify it belongs to the merchant
	if err := database.DB.
		Preload("Event").
		Where("ticket_id = ? AND deleted_at IS NULL", ticketID).
		First(&ticket).Error; err != nil {
		return nil, ErrTicketNotFound
	}

	// Verify event belongs to the merchant
	if ticket.Event == nil || ticket.Event.MerchantID != merchantID {
		return nil, ErrTicketNotFound
	}

	if ticket.Status == "used" {
		return nil, ErrTicketAlreadyUsed
	}

	if ticket.Status != "active" {
		return nil, ErrTicketNotFound
	}

	// Mark as used
	if err := database.DB.Model(&ticket).
		Update("status", "used").Error; err != nil {
		return nil, ErrFailedToUpdate
	}

	slog.Info("Ticket validated",
		"ticket_id", ticketID,
		"merchant_id", merchantID,
	)

	eventResp := toEventResponse(ticket.Event)
	return toTicketResponse(&ticket, eventResp), nil
}

func toTicketResponse(t *models.Ticket, event *dto.EventResponse) *dto.TicketResponse {
	return &dto.TicketResponse{
		TicketID:      t.TicketID,
		TransactionID: t.TransactionID,
		UserID:        t.UserID,
		EventID:       t.EventID,
		Status:        t.Status,
		Event:         event,
		CreatedAt:     t.CreatedAt,
	}
}
