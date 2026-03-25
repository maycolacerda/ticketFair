package dto

import "time"

type AdminLoginRequest struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AdminLoginResponse struct {
	Token     string        `json:"token"`
	ExpiresAt int64         `json:"expires_at"`
	Admin     AdminResponse `json:"admin"`
}

type AdminResponse struct {
	AdminID   string    `json:"admin_id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
}

type AdminCreateUserRequest struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required,password"`
	Username string `json:"username" validate:"required,min=3,max=32,alphanum"`
	Active   *bool  `json:"active"   validate:"omitempty"`
}

type AdminCreateMerchantRequest struct {
	Name        string `json:"name"        validate:"required,min=2,max=100"`
	Email       string `json:"email"       validate:"required,email"`
	Password    string `json:"password"    validate:"required,password"`
	Phone       string `json:"phone"       validate:"required"`
	Description string `json:"description" validate:"omitempty,max=500"`
	Active      *bool  `json:"active"      validate:"omitempty"`
}

type AdminUpdateUserRequest struct {
	Username string `json:"username" validate:"omitempty,min=3,max=32,alphanum"`
	Active   *bool  `json:"active"   validate:"omitempty"`
}

type AdminUpdateMerchantRequest struct {
	Name        string `json:"name"        validate:"omitempty,min=2,max=100"`
	Phone       string `json:"phone"       validate:"omitempty"`
	Description string `json:"description" validate:"omitempty,max=500"`
	Active      *bool  `json:"active"      validate:"omitempty"`
}

type PaginatedAdminUsersResponse struct {
	Data  []UserResponse `json:"data"`
	Page  int            `json:"page"`
	Limit int            `json:"limit"`
	Total int64          `json:"total"`
}

type PaginatedAdminMerchantsResponse struct {
	Data  []MerchantResponse `json:"data"`
	Page  int                `json:"page"`
	Limit int                `json:"limit"`
	Total int64              `json:"total"`
}
