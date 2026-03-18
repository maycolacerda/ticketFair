// services/merchant_rep_service.go
package services

import (
	"errors"
	"strings"

	"github.com/maycolacerda/ticketfair/database"
	"github.com/maycolacerda/ticketfair/dto"
	"github.com/maycolacerda/ticketfair/models"
	"golang.org/x/crypto/bcrypt"
)

func CreateMerchantRep(merchantID string, req dto.CreateMerchantRepRequest) (*dto.MerchantRepResponse, error) {
	// Confirm merchant exists and is active
	var merchant models.Merchant
	if err := database.DB.First(&merchant, "merchant_id = ?", merchantID).Error; err != nil {
		return nil, errors.New("merchant not found")
	}
	if !merchant.Active {
		return nil, errors.New("merchant account is disabled")
	}

	// Email must be globally unique — used as login identifier
	var existing models.MerchantRep
	if err := database.DB.Where("email = ?", req.Email).First(&existing).Error; err == nil {
		return nil, errors.New("email already in use")
	}

	// Hash password before storing
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to process password")
	}

	rep := models.MerchantRep{
		MerchantID: merchantID,
		Name:       strings.TrimSpace(req.Name),
		Email:      strings.ToLower(strings.TrimSpace(req.Email)),
		Phone:      strings.TrimSpace(req.Phone),
		Role:       req.Role,
		Password:   string(hash), // ← was missing
	}

	if err := database.DB.Create(&rep).Error; err != nil {
		return nil, errors.New("failed to create merchant representative")
	}

	return toMerchantRepResponse(&rep), nil
}

func UpdateMerchantRep(merchantID, repID string, req dto.UpdateMerchantRepRequest) (*dto.MerchantRepResponse, error) {
	var rep models.MerchantRep

	// Scope lookup to merchant — prevents updating another merchant's rep
	if err := database.DB.Where("merchant_rep_id = ? AND merchant_id = ?", repID, merchantID).First(&rep).Error; err != nil {
		return nil, errors.New("merchant representative not found")
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = strings.TrimSpace(req.Name)
	}
	if req.Phone != "" {
		updates["phone"] = strings.TrimSpace(req.Phone)
	}
	if req.Role != "" {
		updates["role"] = req.Role
	}

	if len(updates) == 0 {
		return nil, errors.New("no fields to update")
	}

	if err := database.DB.Model(&rep).Updates(updates).Error; err != nil {
		return nil, errors.New("failed to update merchant representative")
	}

	// Re-fetch to return fresh data — Updates() doesn't refresh the local struct
	if err := database.DB.First(&rep, "merchant_rep_id = ?", repID).Error; err != nil {
		return nil, errors.New("failed to fetch updated merchant representative")
	}

	return toMerchantRepResponse(&rep), nil
}

func toMerchantRepResponse(r *models.MerchantRep) *dto.MerchantRepResponse {
	return &dto.MerchantRepResponse{
		MerchantRepID: r.MerchantRepID,
		MerchantID:    r.MerchantID,
		Name:          r.Name,
		Email:         r.Email,
		Phone:         r.Phone,
		Role:          r.Role,
		CreatedAt:     r.CreatedAt,
	}
}
