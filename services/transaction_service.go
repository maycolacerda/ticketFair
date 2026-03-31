// services/transaction_service.go
package services

import (
	"log/slog"
	"strings"

	"github.com/maycolacerda/ticketfair/database"
	"github.com/maycolacerda/ticketfair/dto"
	"github.com/maycolacerda/ticketfair/models"
)

func PurchaseTicket(userID, eventID, ticketTypeID string, quantity int, amount float64) (*dto.TransactionResponse, error) {
	var transactionID string

	err := database.DB.Raw(
		`SELECT purchase_ticket(?, ?, ?, ?, ?)`,
		userID, eventID, ticketTypeID, quantity, amount,
	).Scan(&transactionID).Error

	if err != nil {
		slog.Error("Ticket purchase failed",
			"user_id", userID,
			"event_id", eventID,
			"ticket_type_id", ticketTypeID,
			"error", err.Error(),
		)
		switch {
		case containsError(err, "ticket_type_not_found"):
			return nil, ErrTicketTypeNotFound
		case containsError(err, "ticket_type_sold_out"):
			return nil, ErrTicketTypeSoldOut
		case containsError(err, "ticket_sale_not_started"):
			return nil, ErrTicketSaleNotStarted
		case containsError(err, "ticket_sale_ended"):
			return nil, ErrTicketSaleEnded
		case containsError(err, "ticket_below_minimum"):
			return nil, ErrTicketBelowMinimum
		case containsError(err, "ticket_exceeds_maximum"):
			return nil, ErrTicketExceedsMaximum
		case containsError(err, "event_not_found"):
			return nil, ErrEventNotFound
		default:
			return nil, ErrFailedToCreate
		}
	}

	// Fetch ticket type snapshot for denormalization
	var tt models.TicketType
	_ = database.DB.First(&tt, "ticket_type_id = ?", ticketTypeID).Error

	// Create one ticket per quantity
	for i := 0; i < quantity; i++ {
		ticket := models.Ticket{
			TransactionID:  transactionID,
			UserID:         userID,
			EventID:        eventID,
			TicketTypeID:   ticketTypeID,
			TicketTypeName: tt.Name,
			PricePaidCents: tt.PriceCents,
			Status:         "active",
		}
		if err := database.DB.Create(&ticket).Error; err != nil {
			slog.Error("Failed to create ticket",
				"transaction_id", transactionID,
				"index", i,
				"error", err.Error(),
			)
		}
	}

	slog.Info("Purchase completed",
		"transaction_id", transactionID,
		"user_id", userID,
		"event_id", eventID,
		"ticket_type_id", ticketTypeID,
		"quantity", quantity,
	)

	return GetTransactionByID(transactionID)
}

func RefundTicket(transactionID, userID string) error {
	// Verify ownership before refunding
	var tx models.Transaction
	if err := database.DB.
		Where("transaction_id = ? AND user_id = ?", transactionID, userID).
		First(&tx).Error; err != nil {
		return ErrTransactionNotFound
	}

	// Verify ticket isn't already used
	var ticket models.Ticket
	if err := database.DB.
		Where("transaction_id = ?", transactionID).
		First(&ticket).Error; err == nil {
		if ticket.Status == "used" {
			return ErrNotRefundable
		}
	}

	err := database.DB.Exec(
		`SELECT refund_ticket(?)`, transactionID,
	).Error

	if err != nil {
		slog.Error("Refund failed",
			"transaction_id", transactionID,
			"user_id", userID,
			"error", err.Error(),
		)
		switch {
		case containsError(err, "transaction_not_found"):
			return ErrTransactionNotFound
		case containsError(err, "transaction_not_refundable"):
			return ErrNotRefundable
		default:
			return ErrFailedToUpdate
		}
	}

	// Deactivate ticket
	database.DB.Model(&models.Ticket{}).
		Where("transaction_id = ?", transactionID).
		Update("status", "refunded")

	slog.Info("Refund completed",
		"transaction_id", transactionID,
		"user_id", userID,
	)

	return nil
}

func GetTransactionByID(transactionID string) (*dto.TransactionResponse, error) {
	var tx models.Transaction

	if err := database.DB.
		First(&tx, "transaction_id = ?", transactionID).Error; err != nil {
		return nil, ErrTransactionNotFound
	}

	return toTransactionResponse(&tx), nil
}

func GetUserTransactions(userID string, page, limit int) (*dto.PaginatedTransactionsResponse, error) {
	var transactions []models.Transaction
	var total int64

	offset := (page - 1) * limit

	if err := database.DB.Model(&models.Transaction{}).
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Count(&total).Error; err != nil {
		return nil, ErrFailedToFetch
	}

	if err := database.DB.
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&transactions).Error; err != nil {
		return nil, ErrFailedToFetch
	}

	data := make([]dto.TransactionResponse, len(transactions))
	for i, tx := range transactions {
		data[i] = *toTransactionResponse(&tx)
	}

	return &dto.PaginatedTransactionsResponse{
		Data:  data,
		Page:  page,
		Limit: limit,
		Total: total,
	}, nil
}

func toTransactionResponse(tx *models.Transaction) *dto.TransactionResponse {
	return &dto.TransactionResponse{
		TransactionID: tx.TransactionID,
		UserID:        tx.UserID,
		EventID:       tx.EventID,
		Amount:        tx.Amount,
		Status:        tx.Status,
		CreatedAt:     tx.CreatedAt,
	}
}

func containsError(err error, msg string) bool {
	return err != nil && strings.Contains(err.Error(), msg)
}
