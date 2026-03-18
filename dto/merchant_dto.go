// dto/merchant_dto.go
package dto

import "time"

type CreateMerchantRequest struct {
	Name        string `json:"name"        validate:"required,min=2,max=100"`
	Email       string `json:"email"       validate:"required,email"`
	Phone       string `json:"phone"       validate:"required"`
	Description string `json:"description" validate:"max=500"`
	Password    string `json:"password"    validate:"required,password"`
}

type UpdateMerchantRequest struct {
	Name        string `json:"name"         validate:"omitempty,min=2,max=100"`
	Phone       string `json:"phone"        validate:"omitempty"`
	Description string `json:"description"  validate:"omitempty,max=500"`
}

type MerchantResponse struct {
	MerchantID  string    `json:"merchant_id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Phone       string    `json:"phone"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type MerchantLoginResponse struct {
	Token     string           `json:"token"`
	ExpiresAt int64            `json:"expires_at"`
	Merchant  MerchantResponse `json:"merchant"` // ← named merchant, not user
}
