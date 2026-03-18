package dto

import (
	"time"
)

type CreateUserRequest struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required,password"`
	Username string `json:"username" validate:"required,min=3,max=32,alphanum"`
}

type UserResponse struct {
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

type PaginatedUsersResponse struct {
	Data  []UserResponse `json:"data"`
	Page  int            `json:"page"`
	Limit int            `json:"limit"`
	Total int64          `json:"total"`
}
