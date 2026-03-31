// dto/password_reset_dto.go
package dto

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ForgotPasswordResponse struct {
	Message string `json:"message"`
}

type ResetPasswordRequest struct {
	Email       string `json:"email"        validate:"required,email"`
	Code        string `json:"code"         validate:"required,len=6"`
	NewPassword string `json:"new_password" validate:"required,password"`
}

type ResetPasswordResponse struct {
	Message string `json:"message"`
}
