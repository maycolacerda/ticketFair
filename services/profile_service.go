// services/profile_service.go
package services

import (
	"errors"
	"log/slog"
	"strings"

	"github.com/maycolacerda/ticketfair/database"
	"github.com/maycolacerda/ticketfair/dto"
	"github.com/maycolacerda/ticketfair/models"
	"gorm.io/gorm"
)

func CreateProfile(userID string, req dto.CreateProfileRequest) (*dto.ProfileResponse, error) {
	var user models.User
	if err := database.DB.First(&user, "user_id = ?", userID).Error; err != nil {
		return nil, ErrUserNotFound
	}

	var existing models.Profile
	if err := database.DB.Where("user_id = ?", userID).First(&existing).Error; err == nil {
		return nil, ErrProfileExists
	}

	if err := database.DB.Where("phone_number = ?", req.PhoneNumber).First(&existing).Error; err == nil {
		return nil, ErrPhoneInUse
	}

	// Wrap in transaction — profile and address must both succeed or both fail
	var profile models.Profile
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		profile = models.Profile{
			UserID:      userID,
			FirstName:   strings.TrimSpace(req.FirstName),
			LastName:    strings.TrimSpace(req.LastName),
			PhoneNumber: strings.TrimSpace(req.PhoneNumber),
		}
		if err := tx.Create(&profile).Error; err != nil {
			return errors.New("failed to create profile")
		}

		address := models.Address{
			ProfileID: profile.ProfileID,
			Street:    strings.TrimSpace(req.Address.Street),
			City:      strings.TrimSpace(req.Address.City),
			State:     strings.TrimSpace(req.Address.State),
			Country:   strings.ToUpper(strings.TrimSpace(req.Address.Country)),
			ZipCode:   strings.TrimSpace(req.Address.ZipCode),
		}
		if err := tx.Create(&address).Error; err != nil {
			return errors.New("failed to create address")
		}

		profile.Address = address
		return nil
	})

	if err != nil {
		slog.Error("Failed to create profile", "user_id", userID, "error", err.Error())
		return nil, err
	}

	slog.Info("Profile created", "profile_id", profile.ProfileID, "user_id", userID)
	return toProfileResponse(&profile), nil
}

func GetProfile(userID string) (*dto.ProfileResponse, error) {
	var profile models.Profile

	if err := database.DB.
		Preload("Address").
		Where("user_id = ?", userID).
		First(&profile).Error; err != nil {
		return nil, ErrProfileNotFound // ← was errors.New("profile not found")
	}

	return toProfileResponse(&profile), nil
}

func UpdateProfile(userID string, req dto.UpdateProfileRequest) (*dto.ProfileResponse, error) {
	var profile models.Profile
	if err := database.DB.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		return nil, ErrProfileNotFound
	}

	profileUpdates := map[string]interface{}{}
	addressUpdates := map[string]interface{}{}

	if req.FirstName != "" {
		profileUpdates["first_name"] = strings.TrimSpace(req.FirstName)
	}
	if req.LastName != "" {
		profileUpdates["last_name"] = strings.TrimSpace(req.LastName)
	}
	if req.PhoneNumber != "" {
		var existing models.Profile
		if err := database.DB.
			Where("phone_number = ? AND user_id != ?", req.PhoneNumber, userID).
			First(&existing).Error; err == nil {
			return nil, ErrPhoneInUse
		}
		profileUpdates["phone_number"] = strings.TrimSpace(req.PhoneNumber)
	}

	if req.Address.Street != "" {
		addressUpdates["street"] = strings.TrimSpace(req.Address.Street)
	}
	if req.Address.City != "" {
		addressUpdates["city"] = strings.TrimSpace(req.Address.City)
	}
	if req.Address.State != "" {
		addressUpdates["state"] = strings.TrimSpace(req.Address.State)
	}
	if req.Address.Country != "" {
		addressUpdates["country"] = strings.ToUpper(strings.TrimSpace(req.Address.Country))
	}
	if req.Address.ZipCode != "" {
		addressUpdates["zip_code"] = strings.TrimSpace(req.Address.ZipCode)
	}

	if len(profileUpdates) == 0 && len(addressUpdates) == 0 {
		return nil, ErrNoFieldsToUpdate // ← was errors.New("no fields to update")
	}

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		if len(profileUpdates) > 0 {
			if err := tx.Model(&profile).Updates(profileUpdates).Error; err != nil {
				return ErrFailedToUpdate // ← was errors.New("failed to update profile")
			}
		}
		if len(addressUpdates) > 0 {
			if err := tx.Model(&models.Address{}).
				Where("profile_id = ?", profile.ProfileID).
				Updates(addressUpdates).Error; err != nil {
				return ErrFailedToUpdate // ← was errors.New("failed to update address")
			}
		}
		return nil
	})

	if err != nil {
		slog.Error("Failed to update profile", "user_id", userID, "error", err.Error())
		return nil, err
	}

	if err := database.DB.
		Preload("Address").
		Where("user_id = ?", userID).
		First(&profile).Error; err != nil {
		return nil, ErrFailedToFetch // ← was errors.New("failed to fetch updated profile")
	}

	slog.Info("Profile updated", "profile_id", profile.ProfileID, "user_id", userID)
	return toProfileResponse(&profile), nil
}

func toProfileResponse(p *models.Profile) *dto.ProfileResponse {
	resp := &dto.ProfileResponse{
		ProfileID:     p.ProfileID,
		UserID:        p.UserID,
		FirstName:     p.FirstName,
		LastName:      p.LastName,
		PhoneNumber:   p.PhoneNumber,
		VerifiedEmail: p.VerifiedEmail,
		VerifiedPhone: p.VerifiedPhone,
		CreatedAt:     p.CreatedAt,
	}
	return resp
}
