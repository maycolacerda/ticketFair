package services

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/maycolacerda/ticketfair/database"
	"github.com/maycolacerda/ticketfair/models"
	"golang.org/x/crypto/bcrypt"
)

// NewAuthRequest godoc
//
//	@Summary		Authenticate a user.
//	@Description	Authenticate a user with email and password
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			loginRequest	body	models.LoginRequest	true	"Login request data"
//	@Success		200	{object}	map[string]string
//	@Failure		400	{object}	map[string]string
//	@Failure		404	{object}	map[string]string
//	@Router			/public/auth/login [post]
func NewAuthRequestClient(c *gin.Context) {
	var LoginRequest models.LoginRequest
	var user models.User
	if err := c.ShouldBindJSON(&LoginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body", "details": err.Error()})
		return
	}
	if err := LoginRequest.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	if err := database.DB.First(&user, "email = ?", LoginRequest.Email).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email or password is incorrect"})
		return
	} else {
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(LoginRequest.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Email or password is incorrect"})
			return
		}

		token, err := GenerateClientToken(c, user.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}
		c.Header("Authorization", "Bearer "+token)
		c.JSON(http.StatusOK, gin.H{"message": "Login successful", "user_id": user.UserID})

	}
}

func NewAuthRequestMerchant(c *gin.Context) {
	var LoginRequest models.LoginRequest
	var merchantRep models.MerchantRep
	if err := c.ShouldBindJSON(&LoginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body", "details": err.Error()})
		return
	}
	if err := LoginRequest.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	if err := database.DB.First(&merchantRep, "email = ?", LoginRequest.Email).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email or password is incorrect"})
		return
	} else {
		if err := bcrypt.CompareHashAndPassword([]byte(merchantRep.Password), []byte(LoginRequest.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Email or password is incorrect"})
			return
		}

		token, err := GenerateMerchantToken(c, merchantRep.MerchantRepID, merchantRep.Role)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}
		c.Header("Authorization", "Bearer "+token)
		c.JSON(http.StatusOK, gin.H{"message": "Login successful", "merchant_id": merchantRep.MerchantID, "role": merchantRep.Role, "Merchant Rep ID": merchantRep.MerchantRepID})

	}
}

// Logout godoc
//
//	@Summary		Logout a user.
//	@Description	Logout a user by clearing the Authorization
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	map[string]string
//	@Router			/public/auth/logout [post]
func Logout(c *gin.Context) {
	c.Header("Authorization", "")
	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}
