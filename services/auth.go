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
func NewAuthRequest(c *gin.Context) {
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

		token, err := GenerateToken(c, user.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}
		c.Header("Authorization", "Bearer "+token)
		c.JSON(http.StatusOK, gin.H{"message": "Login successful", "user_id": user.UserID})

	}

}
