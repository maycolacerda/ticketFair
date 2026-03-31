// services/password_reset_service.go
package services

import (
	"log/slog"
	"strings"
	"time"

	"github.com/maycolacerda/ticketfair/database"
	"github.com/maycolacerda/ticketfair/dto"
	"github.com/maycolacerda/ticketfair/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// ─────────────────────────────────────────────
// ForgotPassword
// Looks up the email, generates a 6-digit code,
// saves it and sends the reset email.
// Always returns a generic success message to
// prevent email enumeration attacks.
// ─────────────────────────────────────────────
func ForgotPassword(req dto.ForgotPasswordRequest) (*dto.ForgotPasswordResponse, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))

	generic := &dto.ForgotPasswordResponse{
		Message: "if an account with that email exists, a reset code has been sent",
	}

	// Look up user — return generic response regardless of outcome
	// to prevent email enumeration
	var user models.User
	if err := database.DB.
		Where("email = ? AND active = true AND deleted_at IS NULL", email).
		First(&user).Error; err != nil {
		slog.Warn("Password reset requested for unknown email",
			"email", maskEmail(email),
		)
		return generic, nil
	}

	// Invalidate all previous unused codes for this user
	database.DB.
		Where("user_id = ? AND used_at IS NULL AND deleted_at IS NULL", user.UserID).
		Delete(&models.PasswordReset{})

	// Generate code
	code, err := generateCode()
	if err != nil {
		slog.Error("Failed to generate reset code", "user_id", user.UserID)
		return generic, nil // still generic — don't leak internal errors
	}

	reset := models.PasswordReset{
		UserID:    user.UserID,
		Code:      code,
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}

	if err := database.DB.Create(&reset).Error; err != nil {
		slog.Error("Failed to save reset code", "user_id", user.UserID, "error", err.Error())
		return generic, nil
	}

	// Send email — synchronous so the user gets immediate feedback if SMTP fails
	if err := SendPasswordResetEmail(user.Email, code); err != nil {
		slog.Error("Failed to send reset email",
			"user_id", user.UserID,
			"error", err.Error(),
		)
		// Don't expose SMTP failure — code is saved, user can retry
	}

	slog.Info("Password reset code sent",
		"user_id", user.UserID,
		"email", maskEmail(email),
	)

	return generic, nil
}

// ─────────────────────────────────────────────
// ResetPassword
// Validates the code and sets the new password.
// ─────────────────────────────────────────────
func ResetPassword(req dto.ResetPasswordRequest) (*dto.ResetPasswordResponse, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))

	// Find user
	var user models.User
	if err := database.DB.
		Where("email = ? AND active = true AND deleted_at IS NULL", email).
		First(&user).Error; err != nil {
		return nil, ErrUserNotFound
	}

	// Find the most recent unused code for this user
	var reset models.PasswordReset
	if err := database.DB.
		Where("user_id = ? AND used_at IS NULL AND deleted_at IS NULL", user.UserID).
		Order("created_at DESC").
		First(&reset).Error; err != nil {
		return nil, ErrResetCodeNotFound
	}

	// Check expiry
	if time.Now().After(reset.ExpiresAt) {
		return nil, ErrResetCodeExpired
	}

	// Compare code
	if reset.Code != req.Code {
		slog.Warn("Invalid reset code attempt",
			"user_id", user.UserID,
			"email", maskEmail(email),
		)
		return nil, ErrResetCodeInvalid
	}

	// Hash new password
	hash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, ErrFailedToHash
	}

	// Update password + mark code as used in a single transaction

	txErr := database.DB.Transaction(func(tx *gorm.DB) error {
		now := time.Now()

		if err := tx.Model(&reset).Update("used_at", now).Error; err != nil {
			return err
		}

		if err := tx.Model(&user).Update("password", string(hash)).Error; err != nil {
			return err
		}

		return nil
	})
	if txErr != nil {
		slog.Error("Failed to reset password", "user_id", user.UserID, "error", txErr.Error())
		return nil, ErrFailedToUpdate
	}

	slog.Info("Password reset successful", "user_id", user.UserID)

	return &dto.ResetPasswordResponse{
		Message: "password reset successfully — you can now sign in with your new password",
	}, nil
}
