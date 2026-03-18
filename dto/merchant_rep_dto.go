// dto/merchant_rep_dto.go
package dto

import "time"

type CreateMerchantRepRequest struct {
	Name     string `json:"name"     validate:"required,min=2,max=100"`
	Email    string `json:"email"    validate:"required,email"`
	Phone    string `json:"phone"    validate:"required"`
	Role     string `json:"role"     validate:"required,oneof=admin manager staff"`
	Password string `json:"password" validate:"required,password"`
}

type UpdateMerchantRepRequest struct {
	Name  string `json:"name"  validate:"omitempty,min=2,max=100"`
	Phone string `json:"phone" validate:"omitempty"`
	Role  string `json:"role"  validate:"omitempty,oneof=admin manager staff"`
}

type MerchantRepResponse struct {
	MerchantRepID string    `json:"merchant_rep_id"`
	MerchantID    string    `json:"merchant_id"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	Phone         string    `json:"phone"`
	Role          string    `json:"role"`
	CreatedAt     time.Time `json:"created_at"`
}
type MerchantRepLoginRequest struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type MerchantRepLoginResponse struct {
	Token     string              `json:"token"`
	ExpiresAt int64               `json:"expires_at"`
	Rep       MerchantRepResponse `json:"rep"`
}
