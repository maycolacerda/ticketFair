// services/merchant_rep_auth.go
package services

import (
	"errors"
	"log/slog"
	"strings"

	"github.com/maycolacerda/ticketfair/database"
	"github.com/maycolacerda/ticketfair/dto"
	"github.com/maycolacerda/ticketfair/models"
	"golang.org/x/crypto/bcrypt"
)

func AuthenticateMerchantRep(req dto.MerchantRepLoginRequest) (*dto.MerchantRepLoginResponse, error) {
	var rep models.MerchantRep

	email := strings.ToLower(strings.TrimSpace(req.Email))

	if err := database.DB.
		Preload("Merchant").
		Where("email = ?", email).
		First(&rep).Error; err != nil {
		slog.Warn("Merchant rep login failed — email not found", "email", email)
		return nil, errors.New("invalid credentials")
	}

	if !rep.Active {
		slog.Warn("Merchant rep login failed — account disabled", "rep_id", rep.MerchantRepID)
		return nil, errors.New("account is disabled")
	}

	if !rep.Merchant.Active {
		slog.Warn("Merchant rep login failed — merchant disabled", "merchant_id", rep.MerchantID)
		return nil, errors.New("merchant account is disabled")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(rep.Password), []byte(req.Password)); err != nil {
		slog.Warn("Merchant rep login failed — wrong password", "rep_id", rep.MerchantRepID)
		return nil, errors.New("invalid credentials")
	}

	token, expiresAt, err := GenerateToken(rep.MerchantRepID, rep.Role, rep.MerchantID)
	if err != nil {
		slog.Error("Merchant rep token generation failed", "rep_id", rep.MerchantRepID, "error", err.Error())
		return nil, errors.New("failed to generate token")
	}

	slog.Info("Merchant rep login successful", "rep_id", rep.MerchantRepID, "role", rep.Role, "merchant_id", rep.MerchantID)

	return &dto.MerchantRepLoginResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		Rep: dto.MerchantRepResponse{
			MerchantRepID: rep.MerchantRepID,
			MerchantID:    rep.MerchantID,
			Name:          rep.Name,
			Email:         rep.Email,
			Phone:         rep.Phone,
			Role:          rep.Role,
			CreatedAt:     rep.CreatedAt,
		},
	}, nil
}
