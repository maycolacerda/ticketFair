// controllers/auth.go
package controllers

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/maycolacerda/ticketfair/dto"
	"github.com/maycolacerda/ticketfair/services"
)

// ClientLogin godoc
//
//	@Summary		Client login
//	@Description	Authenticate a client and return a JWT token
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			credentials	body		dto.LoginRequest	true	"Login credentials"
//	@Success		200			{object}	dto.LoginResponse
//	@Failure		400			{object}	map[string]string
//	@Failure		401			{object}	map[string]string
//	@Failure		403			{object}	map[string]string
//	@Router			/public/auth/client/login [post]
func ClientLogin(c *gin.Context) {
	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := validate.Struct(req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": formatValidationErrors(err)})
		return
	}

	resp, err := services.AuthenticateClient(req)
	if err != nil {
		switch err.Error() {
		case "account is disabled":
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		}
		return
	}

	slog.Info("Client login successful", "user_id", resp.User.UserID)
	c.JSON(http.StatusOK, gin.H{"data": resp})
}

// MerchantLogin godoc
//
//	@Summary		Merchant login
//	@Description	Authenticate a merchant and return a JWT token
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			credentials	body		dto.LoginRequest	true	"Login credentials"
//	@Success		200			{object}	dto.MerchantLoginResponse
//	@Failure		400			{object}	map[string]string
//	@Failure		401			{object}	map[string]string
//	@Failure		403			{object}	map[string]string
//	@Router			/public/auth/merchant/login [post]
func MerchantLogin(c *gin.Context) {
	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := validate.Struct(req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": formatValidationErrors(err)})
		return
	}

	resp, err := services.AuthenticateMerchant(req)
	if err != nil {
		switch err.Error() {
		case "account is disabled":
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		}
		return
	}

	slog.Info("Merchant login successful", "merchant_id", resp.Merchant.MerchantID)
	c.JSON(http.StatusOK, gin.H{"data": resp})
}

// MerchantRepLogin godoc
//
//	@Summary		Merchant representative login
//	@Description	Authenticate a merchant rep and return a JWT with role and merchant context
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			credentials	body		dto.MerchantRepLoginRequest		true	"Login credentials"
//	@Success		200			{object}	dto.MerchantRepLoginResponse
//	@Failure		400			{object}	map[string]string
//	@Failure		401			{object}	map[string]string
//	@Failure		403			{object}	map[string]string
//	@Router			/public/auth/rep/login [post]
func MerchantRepLogin(c *gin.Context) {
	var req dto.MerchantRepLoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := validate.Struct(req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": formatValidationErrors(err)})
		return
	}

	resp, err := services.AuthenticateMerchantRep(req)
	if err != nil {
		switch err.Error() {
		case "account is disabled", "merchant account is disabled":
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		}
		return
	}

	slog.Info("Merchant rep login successful", "rep_id", resp.Rep.MerchantRepID, "role", resp.Rep.Role)
	c.JSON(http.StatusOK, gin.H{"data": resp})
}

// Logout godoc
//
//	@Summary		Logout
//	@Description	Invalidate the current session (client-side token removal)
//	@Tags			Auth
//	@Produce		json
//	@Success		200	{object}	dto.LogoutResponse
//	@Router			/public/auth/logout [post]
func Logout(c *gin.Context) {
	slog.Info("User logged out", "user_id", c.GetString("user_id"))
	c.JSON(http.StatusOK, dto.LogoutResponse{Message: "logged out successfully"})
}
