// controllers/password_reset.go
package controllers

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/maycolacerda/ticketfair/dto"
	"github.com/maycolacerda/ticketfair/services"
)

// ForgotPassword godoc
//
//	@Summary		Request a password reset code
//	@Description	Sends a 6-digit reset code to the user's email. Always returns success to prevent email enumeration.
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		dto.ForgotPasswordRequest	true	"Email"
//	@Success		200		{object}	dto.ForgotPasswordResponse
//	@Failure		422		{object}	map[string]interface{}
//	@Failure		429		{object}	map[string]string
//	@Router			/public/auth/password/forgot [post]
func ForgotPassword(c *gin.Context) {
	var req dto.ForgotPasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := validate.Struct(req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": formatValidationErrors(err)})
		return
	}

	resp, err := services.ForgotPassword(req)
	if err != nil {
		// ForgotPassword swallows internal errors and returns generic response
		// so this branch is essentially unreachable, but handle defensively
		slog.Error("ForgotPassword unexpected error", "error", err.Error())
		c.JSON(http.StatusOK, gin.H{
			"message": "if an account with that email exists, a reset code has been sent",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": resp})
}

// ResetPassword godoc
//
//	@Summary		Reset password using a code
//	@Description	Validates the 6-digit code and sets a new password
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		dto.ResetPasswordRequest	true	"Email, code and new password"
//	@Success		200		{object}	dto.ResetPasswordResponse
//	@Failure		400		{object}	map[string]string
//	@Failure		404		{object}	map[string]string
//	@Failure		422		{object}	map[string]interface{}
//	@Failure		429		{object}	map[string]string
//	@Router			/public/auth/password/reset [post]
func ResetPassword(c *gin.Context) {
	var req dto.ResetPasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := validate.Struct(req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": formatValidationErrors(err)})
		return
	}

	resp, err := services.ResetPassword(req)
	if err != nil {
		slog.Warn("ResetPassword failed", "error", err.Error())
		switch {
		case errors.Is(err, services.ErrUserNotFound),
			errors.Is(err, services.ErrResetCodeNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, services.ErrResetCodeExpired),
			errors.Is(err, services.ErrResetCodeInvalid):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to reset password"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": resp})
}
