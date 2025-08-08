package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/maycolacerda/ticketfair/database"
	"github.com/maycolacerda/ticketfair/models"
	"github.com/maycolacerda/ticketfair/services"
)

//	 CreateProfile godoc
//		@Summary		Create a new profile.
//		@Description	Create a new profile with user ID and other details.
//		@Tags			Profiles
//		@Accept			json
//		@Produce		json
//		@Param			profile	body	models.Profile	true	"Profile data"
//		@Success		200	{object}	map[string]string
//		@Failure		400	{object}	map[string]string
//		@Router			/private/profile/new [post]
func CreateProfile(c *gin.Context) {
	var profile models.Profile
	if err := c.ShouldBindJSON(&profile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body", "details": err.Error()})
		return
	}
	profile.UserID, _ = services.ExtractTokenID(c)
	err := profile.Validate()
	if len(err) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid profile data", "details": err})
		return
	} else {
		if err := database.DB.Create(&profile).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create profile"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Profile created successfully", "profile_id": profile.ProfileID})
	}
}

// UpdateProfile godoc
//
//	@Summary		Update an existing profile.
//	@Description	Update an existing profile with new details.
//	@Tags			Profiles
//	@Accept			json
//	@Produce		json
//	@Param			profile	body	models.Profile	true	"Profile data"
//	@Success		200	{object}	map[string]string
//	@Failure		400	{object}	map[string]string
//	@Router			/private/profile/update [post]
func UpdateProfile(c *gin.Context) {
	userID, _ := services.ExtractTokenID(c)
	var profile models.Profile
	if err := c.ShouldBindJSON(&profile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body", "details": err.Error()})
		return
	}

	if err := profile.Validate(); len(err) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid profile data", "details": err})
		return
	}

	if err := database.DB.Model(&models.Profile{}).Where("user_id = ?", userID).Updates(profile).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}
