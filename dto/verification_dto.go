// dto/verification_dto.go
package dto

type SendEmailVerificationRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type SendPhoneVerificationRequest struct {
	Phone string `json:"phone" validate:"required,onlynumbers"`
}

type VerifyEmailRequest struct {
	Email string `json:"email" validate:"required,email"`
	Code  string `json:"code"  validate:"required,len=6"`
}

type VerifyPhoneRequest struct {
	Phone string `json:"phone" validate:"required,onlynumbers"`
	Code  string `json:"code"  validate:"required,len=6"`
}

type VerificationResponse struct {
	Message  string `json:"message"`
	Verified bool   `json:"verified"`
}
