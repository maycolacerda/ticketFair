// services/connection_service.go
package services

import (
	"log/slog"
	"time"

	"github.com/maycolacerda/ticketfair/database"
	"github.com/maycolacerda/ticketfair/dto"
	"github.com/maycolacerda/ticketfair/models"
)

// ─────────────────────────────────────────────
// Send connection request
// ─────────────────────────────────────────────
func SendConnectionRequest(requesterID, addresseeID string) (*dto.ConnectionResponse, error) {
	if requesterID == addresseeID {
		return nil, ErrConnectionSelf
	}

	// Check addressee exists
	var addressee models.User
	if err := database.DB.First(&addressee, "user_id = ? AND active = true", addresseeID).Error; err != nil {
		return nil, ErrUserNotFound
	}

	// Check if connection already exists in either direction
	var existing models.Connection
	if err := database.DB.Where(
		"((requester_id = ? AND addressee_id = ?) OR (requester_id = ? AND addressee_id = ?)) AND deleted_at IS NULL",
		requesterID, addresseeID, addresseeID, requesterID,
	).First(&existing).Error; err == nil {
		return nil, ErrConnectionExists
	}

	conn := models.Connection{
		RequesterID: requesterID,
		AddresseeID: addresseeID,
		Status:      "pending",
	}

	if err := database.DB.Create(&conn).Error; err != nil {
		return nil, ErrFailedToCreate
	}

	slog.Info("Connection request sent",
		"requester_id", requesterID,
		"addressee_id", addresseeID,
	)

	return toConnectionResponse(&conn, addresseeID, &addressee), nil
}

// ─────────────────────────────────────────────
// Respond to a connection request
// ─────────────────────────────────────────────
func RespondToConnection(addresseeID, connectionID, action string) (*dto.ConnectionResponse, error) {
	var conn models.Connection

	// Only the addressee can respond
	if err := database.DB.
		Preload("Requester").
		Where("connection_id = ? AND addressee_id = ? AND deleted_at IS NULL", connectionID, addresseeID).
		First(&conn).Error; err != nil {
		return nil, ErrConnectionNotFound
	}

	if conn.Status != "pending" {
		return nil, ErrConnectionNotPending
	}

	newStatus := "accepted"
	if action == "decline" {
		newStatus = "declined"
	}

	if err := database.DB.Model(&conn).Update("status", newStatus).Error; err != nil {
		return nil, ErrFailedToUpdate
	}

	conn.Status = newStatus

	slog.Info("Connection request responded",
		"connection_id", connectionID,
		"addressee_id", addresseeID,
		"action", action,
	)

	return toConnectionResponse(&conn, conn.RequesterID, conn.Requester), nil
}

// ─────────────────────────────────────────────
// Remove a connection
// ─────────────────────────────────────────────
func RemoveConnection(userID, connectionID string) error {
	var conn models.Connection

	// Either party can remove
	if err := database.DB.Where(
		"connection_id = ? AND (requester_id = ? OR addressee_id = ?) AND deleted_at IS NULL",
		connectionID, userID, userID,
	).First(&conn).Error; err != nil {
		return ErrConnectionNotFound
	}

	if err := database.DB.Delete(&conn).Error; err != nil {
		return ErrFailedToUpdate
	}

	slog.Info("Connection removed", "connection_id", connectionID, "user_id", userID)
	return nil
}

// ─────────────────────────────────────────────
// List connections for a user
// ─────────────────────────────────────────────
func GetConnections(userID string, status string, page, limit int) (*dto.PaginatedConnectionsResponse, error) {
	var connections []models.Connection
	var total int64
	offset := (page - 1) * limit

	query := database.DB.
		Preload("Requester").
		Preload("Addressee").
		Where(
			"(requester_id = ? OR addressee_id = ?) AND deleted_at IS NULL",
			userID, userID,
		)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Model(&models.Connection{}).Count(&total).Error; err != nil {
		return nil, ErrFailedToFetch
	}

	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).
		Find(&connections).Error; err != nil {
		return nil, ErrFailedToFetch
	}

	data := make([]dto.ConnectionResponse, len(connections))
	for i, c := range connections {
		// "the other person" from the perspective of userID
		otherID := c.AddresseeID
		other := c.Addressee
		if c.RequesterID != userID {
			otherID = c.RequesterID
			other = c.Requester
		}
		data[i] = *toConnectionResponse(&c, otherID, other)
	}

	return &dto.PaginatedConnectionsResponse{
		Data:  data,
		Page:  page,
		Limit: limit,
		Total: total,
	}, nil
}

// ─────────────────────────────────────────────
// Get pending incoming requests
// ─────────────────────────────────────────────
func GetPendingRequests(userID string) ([]dto.ConnectionResponse, error) {
	var connections []models.Connection

	if err := database.DB.
		Preload("Requester").
		Where("addressee_id = ? AND status = 'pending' AND deleted_at IS NULL", userID).
		Order("created_at DESC").
		Find(&connections).Error; err != nil {
		return nil, ErrFailedToFetch
	}

	data := make([]dto.ConnectionResponse, len(connections))
	for i, c := range connections {
		data[i] = *toConnectionResponse(&c, c.RequesterID, c.Requester)
	}
	return data, nil
}

// ─────────────────────────────────────────────
// Connection events feed
// ─────────────────────────────────────────────
func GetConnectionEvents(userID string, page, limit int) ([]dto.ConnectionEventResponse, error) {
	// Get all accepted connection partner IDs
	partnerIDs, err := getConnectionPartnerIDs(userID)
	if err != nil || len(partnerIDs) == 0 {
		return []dto.ConnectionEventResponse{}, nil
	}

	// Find events those connections have active tickets for
	type row struct {
		EventID  string
		UserID   string
		Username string
	}

	var rows []row
	if err := database.DB.Raw(`
		SELECT DISTINCT
			t.event_id,
			u.user_id,
			u.username
		FROM tickets t
		JOIN users u ON u.user_id = t.user_id
		WHERE t.user_id IN ?
		  AND t.status = 'active'
		  AND t.deleted_at IS NULL
		  AND u.deleted_at IS NULL
		ORDER BY t.event_id
	`, partnerIDs).Scan(&rows).Error; err != nil {
		return nil, ErrFailedToFetch
	}

	if len(rows) == 0 {
		return []dto.ConnectionEventResponse{}, nil
	}

	// Group by event
	eventMap := make(map[string][]dto.UserSummary)
	for _, r := range rows {
		eventMap[r.EventID] = append(eventMap[r.EventID], dto.UserSummary{
			UserID:   r.UserID,
			Username: r.Username,
		})
	}

	// Fetch event details
	eventIDs := make([]string, 0, len(eventMap))
	for id := range eventMap {
		eventIDs = append(eventIDs, id)
	}

	var events []models.Event
	offset := (page - 1) * limit
	if err := database.DB.
		Where("event_id IN ? AND active = true AND start_time > NOW() AND deleted_at IS NULL", eventIDs).
		Order("start_time ASC").
		Offset(offset).Limit(limit).
		Find(&events).Error; err != nil {
		return nil, ErrFailedToFetch
	}

	result := make([]dto.ConnectionEventResponse, len(events))
	for i, e := range events {
		attendees := eventMap[e.EventID]
		result[i] = dto.ConnectionEventResponse{
			EventID:        e.EventID,
			Name:           e.Name,
			Location:       e.Location,
			StartTime:      e.StartTime,
			ImageURL:       e.ImageURL,
			Attendees:      attendees,
			AttendeesCount: len(attendees),
		}
	}

	return result, nil
}

// ─────────────────────────────────────────────
// Gift a ticket to a connection
// ─────────────────────────────────────────────
func GiftTicket(senderID, ticketID string, req dto.GiftTicketRequest) (*dto.GiftTicketResponse, error) {
	if senderID == req.RecipientID {
		return nil, ErrCannotGiftOwn
	}

	// Load ticket with event
	var ticket models.Ticket
	if err := database.DB.
		Preload("Event").
		Where("ticket_id = ? AND user_id = ? AND deleted_at IS NULL", ticketID, senderID).
		First(&ticket).Error; err != nil {
		return nil, ErrTicketNotFound
	}

	if ticket.Status != "active" {
		return nil, ErrTicketNotActive
	}

	if ticket.IsGift && ticket.GiftedBy != "" {
		return nil, ErrTicketAlreadyGifted
	}

	// Verify they are connected
	if !areConnected(senderID, req.RecipientID) {
		return nil, ErrNotConnected
	}

	// Verify recipient exists
	var recipient models.User
	if err := database.DB.First(&recipient, "user_id = ? AND active = true", req.RecipientID).Error; err != nil {
		return nil, ErrUserNotFound
	}

	now := time.Now()

	// Transfer ticket ownership
	if err := database.DB.Model(&ticket).Updates(map[string]interface{}{
		"user_id":      req.RecipientID,
		"is_gift":      true,
		"gifted_by":    senderID,
		"gifted_at":    now,
		"gift_message": req.Message,
	}).Error; err != nil {
		return nil, ErrFailedToUpdate
	}

	eventName := ""
	if ticket.Event != nil {
		eventName = ticket.Event.Name
	}

	slog.Info("Ticket gifted",
		"ticket_id", ticketID,
		"sender_id", senderID,
		"recipient_id", req.RecipientID,
	)

	return &dto.GiftTicketResponse{
		TicketID:    ticket.TicketID,
		EventName:   eventName,
		RecipientID: req.RecipientID,
		Recipient:   recipient.Username,
		Message:     req.Message,
		GiftedAt:    now,
	}, nil
}

// ─────────────────────────────────────────────
// Helpers
// ─────────────────────────────────────────────

func areConnected(userA, userB string) bool {
	var count int64
	database.DB.Model(&models.Connection{}).Where(
		"((requester_id = ? AND addressee_id = ?) OR (requester_id = ? AND addressee_id = ?)) AND status = 'accepted' AND deleted_at IS NULL",
		userA, userB, userB, userA,
	).Count(&count)
	return count > 0
}

func getConnectionPartnerIDs(userID string) ([]string, error) {
	var connections []models.Connection
	if err := database.DB.
		Where("(requester_id = ? OR addressee_id = ?) AND status = 'accepted' AND deleted_at IS NULL", userID, userID).
		Find(&connections).Error; err != nil {
		return nil, ErrFailedToFetch
	}

	ids := make([]string, 0, len(connections))
	for _, c := range connections {
		if c.RequesterID == userID {
			ids = append(ids, c.AddresseeID)
		} else {
			ids = append(ids, c.RequesterID)
		}
	}
	return ids, nil
}

func toConnectionResponse(c *models.Connection, otherID string, other *models.User) *dto.ConnectionResponse {
	user := dto.UserSummary{UserID: otherID}
	if other != nil {
		user.Username = other.Username
		user.Email = other.Email
	}
	return &dto.ConnectionResponse{
		ConnectionID: c.ConnectionID,
		Status:       c.Status,
		User:         user,
		CreatedAt:    c.CreatedAt,
	}
}
