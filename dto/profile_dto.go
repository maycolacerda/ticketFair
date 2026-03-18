// dto/profile_dto.go
package dto

import "time"

type CreateProfileRequest struct {
	FirstName   string               `json:"first_name"   validate:"required,onlyletters"`
	LastName    string               `json:"last_name"    validate:"required,onlyletters"`
	PhoneNumber string               `json:"phone_number" validate:"required,onlynumbers"`
	Address     CreateAddressRequest `json:"address"      validate:"required"` // ← nested
}

type UpdateProfileRequest struct {
	FirstName   string               `json:"first_name"   validate:"omitempty,onlyletters"`
	LastName    string               `json:"last_name"    validate:"omitempty,onlyletters"`
	PhoneNumber string               `json:"phone_number" validate:"omitempty,onlynumbers"`
	Address     UpdateAddressRequest `json:"address"      validate:"omitempty"` // ← nested
}

type ProfileResponse struct {
	ProfileID     string          `json:"profile_id"`
	UserID        string          `json:"user_id"`
	FirstName     string          `json:"first_name"`
	LastName      string          `json:"last_name"`
	PhoneNumber   string          `json:"phone_number"`
	VerifiedEmail bool            `json:"verified_email"`
	VerifiedPhone bool            `json:"verified_phone"`
	Address       AddressResponse `json:"address"` // ← nested
	CreatedAt     time.Time       `json:"created_at"`
}
