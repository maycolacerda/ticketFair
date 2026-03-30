// controllers/verification.go
package controllers

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/maycolacerda/ticketfair/services"
)

// SendEmailVerification godoc
//
//	@Summary		Send email verification code
//	@Description	Send a 6-digit verification code to the user's email
//	@Tags			Verification
//	@Produce		json
//	@Success		200	{object}	dto.VerificationResponse
//	@Failure		401	{object}	map[string]string
//	@Failure		409	{object}	map[string]string
//	@Router			/private/verify/email/send [post]
func SendEmailVerification(c *gin.Context) {
	userID, err := services.ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	resp, err := services.SendEmailVerification(userID)
	if err != nil {
		slog.Warn("Failed to send email verification", "user_id", userID, "error", err.Error())
		switch {
		case errors.Is(err, services.ErrAlreadyVerified):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		case errors.Is(err, services.ErrUserNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to send verification"})
		}
		return
	}

	slog.Info("Email verification sent", "user_id", userID)
	c.JSON(http.StatusOK, gin.H{"data": resp})
}

// VerifyEmail godoc
//
//	@Summary		Verify email with code
//	@Description	Verify the user's email using the 6-digit code
//	@Tags			Verification
//	@Accept			json
//	@Produce		json
//	@Param			body	body		dto.VerifyEmailRequest	true	"Verification code"
//	@Success		200		{object}	dto.VerificationResponse
//	@Failure		400		{object}	map[string]string
//	@Failure		401		{object}	map[string]string
//	@Failure		404		{object}	map[string]string
//	@Failure		422		{object}	map[string]interface{}
//	@Router			/private/verify/email [post]
func VerifyEmail(c *gin.Context) {
	userID, err := services.ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req struct {
		Code string `json:"code" validate:"required,len=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := validate.Struct(req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": formatValidationErrors(err)})
		return
	}

	resp, err := services.VerifyEmail(userID, req.Code)
	if err != nil {
		slog.Warn("Email verification failed", "user_id", userID, "error", err.Error())
		switch {
		case errors.Is(err, services.ErrVerificationNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, services.ErrVerificationExpired):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, services.ErrVerificationInvalid):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to verify email"})
		}
		return
	}

	slog.Info("Email verified successfully", "user_id", userID)
	c.JSON(http.StatusOK, gin.H{"data": resp})
}

// SendPhoneVerification godoc
//
//	@Summary		Send phone verification code
//	@Description	Send a 6-digit verification code via SMS to the user's phone
//	@Tags			Verification
//	@Produce		json
//	@Success		200	{object}	dto.VerificationResponse
//	@Failure		401	{object}	map[string]string
//	@Failure		404	{object}	map[string]string
//	@Failure		409	{object}	map[string]string
//	@Router			/private/verify/phone/send [post]
func SendPhoneVerification(c *gin.Context) {
	userID, err := services.ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	resp, err := services.SendPhoneVerification(userID)
	if err != nil {
		slog.Warn("Failed to send phone verification", "user_id", userID, "error", err.Error())
		switch {
		case errors.Is(err, services.ErrAlreadyVerified):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		case errors.Is(err, services.ErrProfileNotFound),
			errors.Is(err, services.ErrProfilePhoneNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to send verification"})
		}
		return
	}

	slog.Info("Phone verification sent", "user_id", userID)
	c.JSON(http.StatusOK, gin.H{"data": resp})
}

// VerifyPhone godoc
//
//	@Summary		Verify phone with code
//	@Description	Verify the user's phone using the 6-digit code
//	@Tags			Verification
//	@Accept			json
//	@Produce		json
//	@Param			body	body		dto.VerifyPhoneRequest	true	"Verification code"
//	@Success		200		{object}	dto.VerificationResponse
//	@Failure		400		{object}	map[string]string
//	@Failure		401		{object}	map[string]string
//	@Failure		404		{object}	map[string]string
//	@Failure		422		{object}	map[string]interface{}
//	@Router			/private/verify/phone [post]
func VerifyPhone(c *gin.Context) {
	userID, err := services.ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req struct {
		Code string `json:"code" validate:"required,len=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := validate.Struct(req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": formatValidationErrors(err)})
		return
	}

	resp, err := services.VerifyPhone(userID, req.Code)
	if err != nil {
		slog.Warn("Phone verification failed", "user_id", userID, "error", err.Error())
		switch {
		case errors.Is(err, services.ErrVerificationNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, services.ErrVerificationExpired):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, services.ErrVerificationInvalid):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to verify phone"})
		}
		return
	}

	slog.Info("Phone verified successfully", "user_id", userID)
	c.JSON(http.StatusOK, gin.H{"data": resp})
}
