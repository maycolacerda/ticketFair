// services/merchant_service.go
package services

import (
	"errors"
	"strings"

	"github.com/maycolacerda/ticketfair/database"
	"github.com/maycolacerda/ticketfair/dto"
	"github.com/maycolacerda/ticketfair/models"
	"golang.org/x/crypto/bcrypt"
)

func CreateMerchant(req dto.CreateMerchantRequest) (*dto.MerchantResponse, error) {
	var existing models.Merchant

	email := strings.ToLower(strings.TrimSpace(req.Email))
	if err := database.DB.Where("email = ?", email).First(&existing).Error; err == nil {
		return nil, errors.New("email already in use")
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to process password")
	}

	merchant := models.Merchant{
		Name:        strings.TrimSpace(req.Name),
		Email:       email,
		Password:    string(hash), // ← add
		Phone:       strings.TrimSpace(req.Phone),
		Description: strings.TrimSpace(req.Description),
	}

	if err := database.DB.Create(&merchant).Error; err != nil {
		return nil, errors.New("failed to create merchant")
	}

	return toMerchantResponse(&merchant), nil
}

func UpdateMerchant(merchantID string, req dto.UpdateMerchantRequest) (*dto.MerchantResponse, error) {
	var merchant models.Merchant

	if err := database.DB.First(&merchant, "merchant_id = ?", merchantID).Error; err != nil {
		return nil, errors.New("merchant not found")
	}

	// Only update non-zero fields — never touches email or ID
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

	if len(updates) == 0 {
		return nil, errors.New("no fields to update")
	}

	if err := database.DB.Model(&merchant).Updates(updates).Error; err != nil {
		return nil, errors.New("failed to update merchant")
	}

	return toMerchantResponse(&merchant), nil
}

func GetMerchantByID(merchantID string) (*dto.MerchantResponse, error) {
	var merchant models.Merchant

	if err := database.DB.First(&merchant, "merchant_id = ?", merchantID).Error; err != nil {
		return nil, errors.New("merchant not found")
	}

	return toMerchantResponse(&merchant), nil
}

func toMerchantResponse(m *models.Merchant) *dto.MerchantResponse {
	return &dto.MerchantResponse{
		MerchantID:  m.MerchantID,
		Name:        m.Name,
		Email:       m.Email,
		Phone:       m.Phone,
		Description: m.Description,
		CreatedAt:   m.CreatedAt,
	}
}
