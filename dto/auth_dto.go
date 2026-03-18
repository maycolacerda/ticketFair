package dto

type LoginRequest struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token     string       `json:"token"`
	ExpiresAt int64        `json:"expires_at"` // unix timestamp
	User      UserResponse `json:"user"`
}

type LogoutResponse struct {
	Message string `json:"message"`
}
