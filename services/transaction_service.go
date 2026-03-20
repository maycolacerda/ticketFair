// services/transaction_service.go
package services

import (
	"log/slog"
	"strings"

	"github.com/maycolacerda/ticketfair/database"
	"github.com/maycolacerda/ticketfair/dto"
	"github.com/maycolacerda/ticketfair/models"
)

func PurchaseTicket(userID, eventID string, amount float64) (*dto.TransactionResponse, error) {
	var transactionID string

	err := database.DB.Raw(
		`SELECT purchase_ticket(?, ?, ?)`,
		userID, eventID, amount,
	).Scan(&transactionID).Error

	if err != nil {
		slog.Error("Ticket purchase failed",
			"user_id", userID,
			"event_id", eventID,
			"error", err.Error(),
		)
		switch {
		case containsError(err, "event_not_found"):
			return nil, ErrEventNotFound
		case containsError(err, "event_sold_out"):
			return nil, ErrEventSoldOut
		default:
			return nil, ErrFailedToCreate
		}
	}

	// Create ticket record linked to transaction
	_, ticketErr := CreateTicket(transactionID, userID, eventID)
	if ticketErr != nil {
		// Don't fail the purchase — ticket can be reconciled later
		slog.Error("Failed to create ticket after purchase",
			"transaction_id", transactionID,
			"error", ticketErr.Error(),
		)
	}

	slog.Info("Purchase completed",
		"transaction_id", transactionID,
		"user_id", userID,
		"event_id", eventID,
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
