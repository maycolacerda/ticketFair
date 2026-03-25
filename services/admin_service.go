// services/admin_service.go
package services

import (
	"log/slog"
	"strings"

	"github.com/maycolacerda/ticketfair/database"
	"github.com/maycolacerda/ticketfair/dto"
	"github.com/maycolacerda/ticketfair/models"
	"golang.org/x/crypto/bcrypt"
)

// ─────────────────────────────────────────────
// USERS
// ─────────────────────────────────────────────

func AdminCreateUser(req dto.AdminCreateUserRequest) (*dto.UserResponse, error) {
	var existing models.User

	if err := database.DB.Where("email = ?", req.Email).First(&existing).Error; err == nil {
		return nil, ErrEmailInUse
	}
	if err := database.DB.Where("username = ?", req.Username).First(&existing).Error; err == nil {
		return nil, ErrUsernameInUse
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, ErrFailedToHash
	}

	active := true
	if req.Active != nil {
		active = *req.Active
	}

	user := models.User{
		Email:    strings.ToLower(strings.TrimSpace(req.Email)),
		Username: strings.TrimSpace(req.Username),
		Password: string(hash),
		Active:   active,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		return nil, ErrFailedToCreate
	}

	slog.Info("Admin created user", "user_id", user.UserID)
	return &dto.UserResponse{
		UserID:    user.UserID,
		Email:     user.Email,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
	}, nil
}

func AdminUpdateUser(userID string, req dto.AdminUpdateUserRequest) (*dto.UserResponse, error) {
	var user models.User

	if err := database.DB.First(&user, "user_id = ?", userID).Error; err != nil {
		return nil, ErrUserNotFound
	}

	updates := map[string]interface{}{}

	if req.Username != "" {
		var existing models.User
		if err := database.DB.
			Where("username = ? AND user_id != ?", req.Username, userID).
			First(&existing).Error; err == nil {
			return nil, ErrUsernameInUse
		}
		updates["username"] = strings.TrimSpace(req.Username)
	}
	if req.Active != nil {
		updates["active"] = *req.Active
	}

	if len(updates) == 0 {
		return nil, ErrNoFieldsToUpdate
	}

	if err := database.DB.Model(&user).Updates(updates).Error; err != nil {
		return nil, ErrFailedToUpdate
	}

	if err := database.DB.First(&user, "user_id = ?", userID).Error; err != nil {
		return nil, ErrFailedToFetch
	}

	slog.Info("Admin updated user", "user_id", userID)
	return &dto.UserResponse{
		UserID:    user.UserID,
		Email:     user.Email,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
	}, nil
}

func AdminDeactivateUser(userID string) error {
	var user models.User
	if err := database.DB.First(&user, "user_id = ?", userID).Error; err != nil {
		return ErrUserNotFound
	}

	if err := database.DB.Model(&user).Update("active", false).Error; err != nil {
		return ErrFailedToUpdate
	}

	slog.Info("Admin deactivated user", "user_id", userID)
	return nil
}

func AdminActivateUser(userID string) error {
	var user models.User
	if err := database.DB.First(&user, "user_id = ?", userID).Error; err != nil {
		return ErrUserNotFound
	}

	if err := database.DB.Model(&user).Update("active", true).Error; err != nil {
		return ErrFailedToUpdate
	}

	slog.Info("Admin activated user", "user_id", userID)
	return nil
}

func AdminGetAllUsers(page, limit int) (*dto.PaginatedAdminUsersResponse, error) {
	var users []models.User
	var total int64

	offset := (page - 1) * limit

	// Admin sees ALL users including inactive — no active filter
	if err := database.DB.Model(&models.User{}).
		Where("deleted_at IS NULL").
		Count(&total).Error; err != nil {
		return nil, ErrFailedToFetch
	}

	if err := database.DB.
		Where("deleted_at IS NULL").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&users).Error; err != nil {
		return nil, ErrFailedToFetch
	}

	data := make([]dto.UserResponse, len(users))
	for i, u := range users {
		data[i] = dto.UserResponse{
			UserID:    u.UserID,
			Email:     u.Email,
			Username:  u.Username,
			CreatedAt: u.CreatedAt,
		}
	}

	return &dto.PaginatedAdminUsersResponse{
		Data:  data,
		Page:  page,
		Limit: limit,
		Total: total,
	}, nil
}

// ─────────────────────────────────────────────
// MERCHANTS
// ─────────────────────────────────────────────

func AdminCreateMerchant(req dto.AdminCreateMerchantRequest) (*dto.MerchantResponse, error) {
	var existing models.Merchant
	email := strings.ToLower(strings.TrimSpace(req.Email))

	if err := database.DB.Where("email = ?", email).First(&existing).Error; err == nil {
		return nil, ErrEmailInUse
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, ErrFailedToHash
	}

	active := true
	if req.Active != nil {
		active = *req.Active
	}

	merchant := models.Merchant{
		Name:        strings.TrimSpace(req.Name),
		Email:       email,
		Password:    string(hash),
		Phone:       strings.TrimSpace(req.Phone),
		Description: strings.TrimSpace(req.Description),
		Active:      active,
	}

	if err := database.DB.Create(&merchant).Error; err != nil {
		return nil, ErrFailedToCreate
	}

	slog.Info("Admin created merchant", "merchant_id", merchant.MerchantID)
	return toMerchantResponse(&merchant), nil
}

func AdminUpdateMerchant(merchantID string, req dto.AdminUpdateMerchantRequest) (*dto.MerchantResponse, error) {
	var merchant models.Merchant

	if err := database.DB.First(&merchant, "merchant_id = ?", merchantID).Error; err != nil {
		return nil, ErrMerchantNotFound
	}

	updates := map[string]interface{}{}

	if req.Name != "" {
		updates["name"] = strings.TrimSpace(req.Name)
	}
	if req.Phone != "" {
		updates["phone"] = strings.TrimSpace(req.Phone)
	}
	if req.Description != "" {
		updates["description"] = strings.TrimSpace(req.Description)
	}
	if req.Active != nil {
		updates["active"] = *req.Active
	}

	if len(updates) == 0 {
		return nil, ErrNoFieldsToUpdate
	}

	if err := database.DB.Model(&merchant).Updates(updates).Error; err != nil {
		return nil, ErrFailedToUpdate
	}

	if err := database.DB.First(&merchant, "merchant_id = ?", merchantID).Error; err != nil {
		return nil, ErrFailedToFetch
	}

	slog.Info("Admin updated merchant", "merchant_id", merchantID)
	return toMerchantResponse(&merchant), nil
}

func AdminDeactivateMerchant(merchantID string) error {
	var merchant models.Merchant
	if err := database.DB.First(&merchant, "merchant_id = ?", merchantID).Error; err != nil {
		return ErrMerchantNotFound
	}

	if err := database.DB.Model(&merchant).Update("active", false).Error; err != nil {
		return ErrFailedToUpdate
	}

	slog.Info("Admin deactivated merchant", "merchant_id", merchantID)
	return nil
}

func AdminActivateMerchant(merchantID string) error {
	var merchant models.Merchant
	if err := database.DB.First(&merchant, "merchant_id = ?", merchantID).Error; err != nil {
		return ErrMerchantNotFound
	}

	if err := database.DB.Model(&merchant).Update("active", true).Error; err != nil {
		return ErrFailedToUpdate
	}

	slog.Info("Admin activated merchant", "merchant_id", merchantID)
	return nil
}

func AdminGetAllMerchants(page, limit int) (*dto.PaginatedAdminMerchantsResponse, error) {
	var merchants []models.Merchant
	var total int64

	offset := (page - 1) * limit

	if err := database.DB.Model(&models.Merchant{}).
		Where("deleted_at IS NULL").
		Count(&total).Error; err != nil {
		return nil, ErrFailedToFetch
	}

	if err := database.DB.
		Where("deleted_at IS NULL").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&merchants).Error; err != nil {
		return nil, ErrFailedToFetch
	}

	data := make([]dto.MerchantResponse, len(merchants))
	for i, m := range merchants {
		data[i] = *toMerchantResponse(&m)
	}

	return &dto.PaginatedAdminMerchantsResponse{
		Data:  data,
		Page:  page,
		Limit: limit,
		Total: total,
	}, nil
}
