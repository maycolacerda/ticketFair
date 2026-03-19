// services/auth.go — service layer only, no HTTP, no gin
package services

import (
	"log/slog"
	"strings"

	"github.com/maycolacerda/ticketfair/database"
	"github.com/maycolacerda/ticketfair/dto"
	"github.com/maycolacerda/ticketfair/models"
	"golang.org/x/crypto/bcrypt"
)

func AuthenticateClient(req dto.LoginRequest) (*dto.LoginResponse, error) {
	var user models.User
	email := strings.ToLower(strings.TrimSpace(req.Email))

	if err := database.DB.Where("email = ?", email).First(&user).Error; err != nil {
		slog.Warn("Client login failed — email not found", "email", email)
		return nil, ErrInvalidCredentials
	}

	if !user.Active {
		slog.Warn("Client login failed — account disabled", "user_id", user.UserID)
		return nil, ErrAccountDisabled
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		slog.Warn("Client login failed — wrong password", "user_id", user.UserID)
		return nil, ErrInvalidCredentials
	}

	token, expiresAt, err := GenerateToken(user.UserID, RoleClient, "")
	if err != nil {
		slog.Error("Client token generation failed", "user_id", user.UserID, "error", err.Error())
		return nil, ErrFailedToGenerateToken
	}

	slog.Info("Client login successful", "user_id", user.UserID)
	return &dto.LoginResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		User: dto.UserResponse{
			UserID:    user.UserID,
			Email:     user.Email,
			Username:  user.Username,
			CreatedAt: user.CreatedAt,
		},
	}, nil
}

func AuthenticateMerchant(req dto.LoginRequest) (*dto.MerchantLoginResponse, error) {
	var merchant models.Merchant
	email := strings.ToLower(strings.TrimSpace(req.Email))

	if err := database.DB.Where("email = ?", email).First(&merchant).Error; err != nil {
		slog.Warn("Merchant login failed — email not found", "email", email)
		return nil, ErrInvalidCredentials
	}

	if !merchant.Active {
		slog.Warn("Merchant login failed — account disabled", "merchant_id", merchant.MerchantID)
		return nil, ErrAccountDisabled
	}

	if err := bcrypt.CompareHashAndPassword([]byte(merchant.Password), []byte(req.Password)); err != nil {
		slog.Warn("Merchant login failed — wrong password", "merchant_id", merchant.MerchantID)
		return nil, ErrInvalidCredentials
	}

	token, expiresAt, err := GenerateToken(merchant.MerchantID, RoleMerchant, merchant.MerchantID)
	if err != nil {
		slog.Error("Merchant token generation failed", "merchant_id", merchant.MerchantID, "error", err.Error())
		return nil, ErrFailedToGenerateToken
	}

	slog.Info("Merchant login successful", "merchant_id", merchant.MerchantID)
	return &dto.MerchantLoginResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		Merchant: dto.MerchantResponse{
			MerchantID:  merchant.MerchantID,
			Name:        merchant.Name,
			Email:       merchant.Email,
			Phone:       merchant.Phone,
			Description: merchant.Description,
			CreatedAt:   merchant.CreatedAt,
		},
	}, nil
}
