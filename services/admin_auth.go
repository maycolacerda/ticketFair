package services

import (
	"log/slog"
	"os"
	"strings"

	"github.com/maycolacerda/ticketfair/database"
	"github.com/maycolacerda/ticketfair/dto"
	"github.com/maycolacerda/ticketfair/models"
)

func AuthenticateAdmin(req dto.AdminLoginRequest) (*dto.AdminLoginResponse, error) {
	var admin models.Admin
	email := strings.ToLower(strings.TrimSpace(req.Email))

	if err := database.DB.Where("email = ?", email).First(&admin).Error; err != nil {
		slog.Warn("Admin login failed — email not found", "email", email)
		return nil, ErrInvalidCredentials
	}

	if !admin.Active {
		slog.Warn("Admin login failed — account disabled", "admin_id", admin.AdminID)
		return nil, ErrAdminDisabled
	}

	if err := admin.Password == os.Getenv("ADMIN_PASSWORD"); !err {
		slog.Warn("Admin login failed — wrong password", "admin_id", admin.AdminID)
		return nil, ErrInvalidCredentials
	}

	token, expiresAt, err := GenerateToken(admin.AdminID, RoleAdmin, "")
	if err != nil {
		slog.Error("Admin token generation failed", "admin_id", admin.AdminID, "error", err.Error())
		return nil, ErrFailedToGenerateToken
	}

	slog.Info("Admin login successful", "admin_id", admin.AdminID)

	return &dto.AdminLoginResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		Admin: dto.AdminResponse{
			AdminID:   admin.AdminID,
			Name:      admin.Name,
			Email:     admin.Email,
			Active:    admin.Active,
			CreatedAt: admin.CreatedAt,
		},
	}, nil
}
