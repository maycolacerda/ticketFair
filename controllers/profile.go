// controllers/profile.go
package controllers

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/maycolacerda/ticketfair/dto"
	"github.com/maycolacerda/ticketfair/services"
)

// CreateProfile godoc
//
//	@Summary		Create a user profile
//	@Description	Create a profile for the authenticated user
//	@Tags			Profile
//	@Accept			json
//	@Produce		json
//	@Param			profile	body		dto.CreateProfileRequest	true	"Profile data"
//	@Success		201		{object}	dto.ProfileResponse
//	@Failure		400		{object}	map[string]string
//	@Failure		401		{object}	map[string]string
//	@Failure		409		{object}	map[string]string
//	@Failure		422		{object}	map[string]interface{}
//	@Router			/private/profile [post]
func CreateProfile(c *gin.Context) {
	userID, err := services.ExtractTokenID(c)
	if err != nil {
		slog.Warn("Unauthorized profile creation attempt")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req dto.CreateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Warn("Invalid request body", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := validate.Struct(req); err != nil {
		errs := formatValidationErrors(err)
		slog.Warn("Profile validation failed", "errors", errs)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": errs})
		return
	}

	profile, err := services.CreateProfile(userID, req)
	if err != nil {
		slog.Warn("Profile creation failed", "user_id", userID, "error", err.Error())
		switch err.Error() {
		case "user not found":
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case "profile already exists", "phone number already in use":
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create profile"})
		}
		return
	}

	slog.Info("Profile created", "profile_id", profile.ProfileID, "user_id", userID)
	c.JSON(http.StatusCreated, gin.H{"data": profile})
}

// UpdateProfile godoc
//
//	@Summary		Update a user profile
//	@Description	Update the authenticated user's profile
//	@Tags			Profile
//	@Accept			json
//	@Produce		json
//	@Param			profile	body		dto.UpdateProfileRequest	true	"Updated profile data"
//	@Success		200		{object}	dto.ProfileResponse
//	@Failure		400		{object}	map[string]string
//	@Failure		401		{object}	map[string]string
//	@Failure		404		{object}	map[string]string
//	@Failure		422		{object}	map[string]interface{}
//	@Router			/private/profile [put]
func UpdateProfile(c *gin.Context) {
	userID, err := services.ExtractTokenID(c)
	if err != nil {
		slog.Warn("Unauthorized profile update attempt")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Warn("Invalid request body", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := validate.Struct(req); err != nil {
		errs := formatValidationErrors(err)
		slog.Warn("Profile update validation failed", "errors", errs)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": errs})
		return
	}

	profile, err := services.UpdateProfile(userID, req)
	if err != nil {
		slog.Warn("Profile update failed", "user_id", userID, "error", err.Error())
		switch err.Error() {
		case "profile not found":
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case "phone number already in use":
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		case "no fields to update":
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update profile"})
		}
		return
	}

	slog.Info("Profile updated", "profile_id", profile.ProfileID, "user_id", userID)
	c.JSON(http.StatusOK, gin.H{"data": profile})
}

// GetProfile godoc
//
//	@Summary		Get the current user's profile
//	@Description	Retrieve the authenticated user's profile
//	@Tags			Profile
//	@Produce		json
//	@Success		200	{object}	dto.ProfileResponse
//	@Failure		401	{object}	map[string]string
//	@Failure		404	{object}	map[string]string
//	@Router			/private/profile [get]
func GetProfile(c *gin.Context) {
	userID, err := services.ExtractTokenID(c)
	if err != nil {
		slog.Warn("Unauthorized profile fetch attempt")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	profile, err := services.GetProfile(userID)
	if err != nil {
		slog.Warn("Profile not found", "user_id", userID)
		c.JSON(http.StatusNotFound, gin.H{"error": "profile not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": profile})
}
