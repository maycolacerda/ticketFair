// services/verification_service.go
package services

import (
	"crypto/rand"
	"fmt"
	"log/slog"
	"time"

	"github.com/maycolacerda/ticketfair/database"
	"github.com/maycolacerda/ticketfair/dto"
	"github.com/maycolacerda/ticketfair/models"
)

// generateCode generates a 6-digit numeric code
func generateCode() (string, error) {
	b := make([]byte, 3)
	if _, err := rand.Read(b); err != nil {
		return "", ErrFailedToCreate
	}
	// Convert to 6-digit number
	code := fmt.Sprintf("%06d", int(b[0])<<16|int(b[1])<<8|int(b[2]))
	return code[len(code)-6:], nil
}

// ─────────────────────────────────────────────
// EMAIL VERIFICATION
// ─────────────────────────────────────────────

func SendEmailVerification(userID string) (*dto.VerificationResponse, error) {
	var user models.User
	if err := database.DB.First(&user, "user_id = ?", userID).Error; err != nil {
		return nil, ErrUserNotFound
	}

	if isEmailVerified(userID) {
		return nil, ErrAlreadyVerified
	}

	database.DB.
		Where("user_id = ? AND type = ? AND used_at IS NULL", userID, "email").
		Delete(&models.Verification{})

	code, err := generateCode()
	if err != nil {
		return nil, err
	}

	verification := models.Verification{
		UserID:    userID,
		Type:      "email",
		Code:      code,
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}

	if err := database.DB.Create(&verification).Error; err != nil {
		return nil, ErrFailedToCreate
	}

	// Send email — non-blocking, log error but don't fail the request
	if err := SendVerificationEmail(user.Email, code); err != nil {
		slog.Error("Failed to send verification email",
			"user_id", userID,
			"error", err.Error(),
		)
	}

	slog.Info("Email verification sent", "user_id", userID, "email", maskEmail(user.Email))

	return &dto.VerificationResponse{
		Message:  "verification code sent to " + maskEmail(user.Email),
		Verified: false,
	}, nil
}
func VerifyEmail(userID, code string) (*dto.VerificationResponse, error) {
	var verification models.Verification

	if err := database.DB.
		Where("user_id = ? AND type = ? AND used_at IS NULL AND deleted_at IS NULL", userID, "email").
		Order("created_at DESC").
		First(&verification).Error; err != nil {
		return nil, ErrVerificationNotFound
	}

	if time.Now().After(verification.ExpiresAt) {
		return nil, ErrVerificationExpired
	}

	if verification.Code != code {
		return nil, ErrVerificationInvalid
	}

	// Mark code as used
	now := time.Now()
	if err := database.DB.Model(&verification).Update("used_at", now).Error; err != nil {
		return nil, ErrFailedToUpdate
	}

	// Mark email as verified on profile
	if err := database.DB.Model(&models.Profile{}).
		Where("user_id = ?", userID).
		Update("verified_email", true).Error; err != nil {
		return nil, ErrFailedToUpdate
	}

	slog.Info("Email verified", "user_id", userID)

	return &dto.VerificationResponse{
		Message:  "email verified successfully",
		Verified: true,
	}, nil
}

// ─────────────────────────────────────────────
// PHONE VERIFICATION
// ─────────────────────────────────────────────

func SendPhoneVerification(userID string) (*dto.VerificationResponse, error) {
	// Find profile to get phone number
	var profile models.Profile
	if err := database.DB.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		return nil, ErrProfileNotFound
	}

	if profile.PhoneNumber == "" {
		return nil, ErrProfilePhoneNotFound
	}

	if profile.VerifiedPhone {
		return nil, ErrAlreadyVerified
	}

	// Invalidate previous unused codes
	database.DB.
		Where("user_id = ? AND type = ? AND used_at IS NULL", userID, "phone").
		Delete(&models.Verification{})

	code, err := generateCode()
	if err != nil {
		return nil, err
	}

	verification := models.Verification{
		UserID:    userID,
		Type:      "phone",
		Code:      code,
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}

	if err := database.DB.Create(&verification).Error; err != nil {
		return nil, ErrFailedToCreate
	}

	// TODO: send SMS via Twilio/AWS SNS
	// For now log the code — replace with real SMS sending
	slog.Info("Phone verification code generated",
		"user_id", userID,
		"phone", maskPhone(profile.PhoneNumber),
		"code", code, // ← remove in production
		"expires_at", verification.ExpiresAt,
	)

	return &dto.VerificationResponse{
		Message:  "verification code sent to " + maskPhone(profile.PhoneNumber),
		Verified: false,
	}, nil
}

func VerifyPhone(userID, code string) (*dto.VerificationResponse, error) {
	var verification models.Verification

	if err := database.DB.
		Where("user_id = ? AND type = ? AND used_at IS NULL AND deleted_at IS NULL", userID, "phone").
		Order("created_at DESC").
		First(&verification).Error; err != nil {
		return nil, ErrVerificationNotFound
	}

	if time.Now().After(verification.ExpiresAt) {
		return nil, ErrVerificationExpired
	}

	if verification.Code != code {
		return nil, ErrVerificationInvalid
	}

	now := time.Now()
	if err := database.DB.Model(&verification).Update("used_at", now).Error; err != nil {
		return nil, ErrFailedToUpdate
	}

	// Mark phone as verified on profile
	if err := database.DB.Model(&models.Profile{}).
		Where("user_id = ?", userID).
		Update("verified_phone", true).Error; err != nil {
		return nil, ErrFailedToUpdate
	}

	slog.Info("Phone verified", "user_id", userID)

	return &dto.VerificationResponse{
		Message:  "phone verified successfully",
		Verified: true,
	}, nil
}

// ─────────────────────────────────────────────
// Helpers
// ─────────────────────────────────────────────

func isEmailVerified(userID string) bool {
	var profile models.Profile
	if err := database.DB.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		return false
	}
	return profile.VerifiedEmail
}

// maskEmail: test@example.com → t***@example.com
func maskEmail(email string) string {
	for i, c := range email {
		if c == '@' {
			return string(email[0]) + "***" + email[i:]
		}
	}
	return "***"
}

// maskPhone: 44999999999 → 44*******99
func maskPhone(phone string) string {
	if len(phone) < 4 {
		return "***"
	}
	return phone[:2] + "****" + phone[len(phone)-2:]
}
