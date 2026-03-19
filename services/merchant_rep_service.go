// services/merchant_rep_service.go
package services

import (
	"strings"

	"github.com/maycolacerda/ticketfair/database"
	"github.com/maycolacerda/ticketfair/dto"
	"github.com/maycolacerda/ticketfair/models"
	"golang.org/x/crypto/bcrypt"
)

func CreateMerchantRep(merchantID string, req dto.CreateMerchantRepRequest) (*dto.MerchantRepResponse, error) {
	var merchant models.Merchant
	if err := database.DB.First(&merchant, "merchant_id = ?", merchantID).Error; err != nil {
		return nil, ErrMerchantNotFound
	}
	if !merchant.Active {
		return nil, ErrMerchantDisabled
	}

	var existing models.MerchantRep
	if err := database.DB.Where("email = ?", req.Email).First(&existing).Error; err == nil {
		return nil, ErrEmailInUse
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, ErrFailedToHash
	}

	rep := models.MerchantRep{
		MerchantID: merchantID,
		Name:       strings.TrimSpace(req.Name),
		Email:      strings.ToLower(strings.TrimSpace(req.Email)),
		Phone:      strings.TrimSpace(req.Phone),
		Role:       req.Role,
		Password:   string(hash),
	}

	if err := database.DB.Create(&rep).Error; err != nil {
		return nil, ErrFailedToCreate
	}

	return toMerchantRepResponse(&rep), nil
}

func UpdateMerchantRep(merchantID, repID string, req dto.UpdateMerchantRepRequest) (*dto.MerchantRepResponse, error) {
	var rep models.MerchantRep

	if err := database.DB.Where("merchant_rep_id = ? AND merchant_id = ?", repID, merchantID).First(&rep).Error; err != nil {
		return nil, ErrRepNotFound
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
		return nil, ErrNoFieldsToUpdate
	}

	if err := database.DB.Model(&rep).Updates(updates).Error; err != nil {
		return nil, ErrFailedToUpdate
	}

	if err := database.DB.First(&rep, "merchant_rep_id = ?", repID).Error; err != nil {
		return nil, ErrFailedToFetch
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
