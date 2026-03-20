// services/event_service.go
package services

import (
	"log/slog"
	"strings"
	"time"

	"github.com/maycolacerda/ticketfair/database"
	"github.com/maycolacerda/ticketfair/dto"
	"github.com/maycolacerda/ticketfair/models"
)

func CreateEvent(merchantID string, req dto.CreateEventRequest) (*dto.EventResponse, error) {
	var merchant models.Merchant
	if err := database.DB.First(&merchant, "merchant_id = ?", merchantID).Error; err != nil {
		return nil, ErrMerchantNotFound
	}
	if !merchant.Active {
		return nil, ErrMerchantDisabled
	}
	if !req.EndTime.After(req.StartTime) {
		return nil, ErrInvalidTimeRange
	}
	if req.StartTime.Before(time.Now()) {
		return nil, ErrStartTimeInPast
	}

	event := models.Event{
		MerchantID:  merchantID,
		Name:        strings.TrimSpace(req.Name),
		Description: strings.TrimSpace(req.Description),
		Location:    strings.TrimSpace(req.Location),
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		Capacity:    req.Capacity,
	}

	if err := database.DB.Create(&event).Error; err != nil {
		return nil, ErrFailedToCreate
	}

	return toEventResponse(&event), nil
}

func UpdateEvent(merchantID, eventID string, req dto.UpdateEventRequest) (*dto.EventResponse, error) {
	var event models.Event

	if err := database.DB.Where("event_id = ? AND merchant_id = ?", eventID, merchantID).First(&event).Error; err != nil {
		return nil, ErrEventNotFound
	}

	updates := map[string]interface{}{}

	if req.Name != "" {
		updates["name"] = strings.TrimSpace(req.Name)
	}
	if req.Description != "" {
		updates["description"] = strings.TrimSpace(req.Description)
	}
	if req.Location != "" {
		updates["location"] = strings.TrimSpace(req.Location)
	}
	if !req.StartTime.IsZero() {
		if req.StartTime.Before(time.Now()) {
			return nil, ErrStartTimeInPast
		}
		updates["start_time"] = req.StartTime
	}
	if !req.EndTime.IsZero() {
		startTime, _ := updates["start_time"].(time.Time)
		if startTime.IsZero() {
			startTime = event.StartTime
		}
		if !req.EndTime.After(startTime) {
			return nil, ErrInvalidTimeRange
		}
		updates["end_time"] = req.EndTime
	}
	if req.Capacity > 0 {
		updates["capacity"] = req.Capacity
	}
	if req.Active != nil {
		updates["active"] = *req.Active
	}

	if len(updates) == 0 {
		return nil, ErrNoFieldsToUpdate
	}

	if err := database.DB.Model(&event).Updates(updates).Error; err != nil {
		slog.Error("Failed to update event", "event_id", eventID, "error", err.Error())
		return nil, ErrFailedToUpdate
	}

	if err := database.DB.First(&event, "event_id = ?", eventID).Error; err != nil {
		return nil, ErrFailedToFetch
	}

	slog.Info("Event updated", "event_id", eventID, "merchant_id", merchantID)
	return toEventResponse(&event), nil
}

func GetEventByID(eventID string) (*dto.EventResponse, error) {
	var event models.Event

	if err := database.DB.First(&event, "event_id = ?", eventID).Error; err != nil {
		return nil, ErrEventNotFound
	}

	return toEventResponse(&event), nil
}

func GetEvents(page, limit int) (*dto.PaginatedEventsResponse, error) {
	var events []models.Event
	var total int64

	offset := (page - 1) * limit

	if err := database.DB.Model(&models.Event{}).
		Where("active = ? AND start_time > ?", true, time.Now()).
		Count(&total).Error; err != nil {
		return nil, ErrFailedToFetch // ← was errors.New("failed to count events")
	}

	if err := database.DB.
		Where("active = ? AND start_time > ?", true, time.Now()).
		Order("start_time ASC").
		Offset(offset).
		Limit(limit).
		Find(&events).Error; err != nil {
		return nil, ErrFailedToFetch // ← was errors.New("failed to fetch events")
	}

	data := make([]dto.EventResponse, len(events))
	for i, e := range events {
		data[i] = *toEventResponse(&e)
	}

	return &dto.PaginatedEventsResponse{
		Data:  data,
		Page:  page,
		Limit: limit,
		Total: total,
	}, nil
}

func toEventResponse(e *models.Event) *dto.EventResponse {
	return &dto.EventResponse{
		EventID:     e.EventID,
		MerchantID:  e.MerchantID,
		Name:        e.Name,
		Description: e.Description,
		Location:    e.Location,
		StartTime:   e.StartTime,
		EndTime:     e.EndTime,
		Capacity:    e.Capacity,
		Active:      e.Active,
		CreatedAt:   e.CreatedAt,
	}
}
