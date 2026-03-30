// controllers/admin.go
package controllers

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/maycolacerda/ticketfair/dto"
	"github.com/maycolacerda/ticketfair/services"
)

// ─────────────────────────────────────────────
// AUTH
// ─────────────────────────────────────────────

// AdminLogin godoc
//
//	@Summary		Admin login
//	@Description	Authenticate an admin and return a JWT token
//	@Tags			Admin
//	@Accept			json
//	@Produce		json
//	@Param			credentials	body		dto.AdminLoginRequest	true	"Admin credentials"
//	@Success		200			{object}	dto.AdminLoginResponse
//	@Failure		400			{object}	map[string]string
//	@Failure		401			{object}	map[string]string
//	@Failure		403			{object}	map[string]string
//	@Router			/admin/auth/login [post]
func AdminLogin(c *gin.Context) {
	var req dto.AdminLoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := validate.Struct(req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": formatValidationErrors(err)})
		return
	}

	resp, err := services.AuthenticateAdmin(req)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrAdminDisabled):
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		}
		return
	}

	slog.Info("Admin login successful", "admin_id", resp.Admin.AdminID)
	c.JSON(http.StatusOK, gin.H{"data": resp})
}

// ─────────────────────────────────────────────
// USERS
// ─────────────────────────────────────────────

// AdminListUsers godoc
//
//	@Summary		List all users
//	@Description	Admin — retrieve all users including inactive
//	@Tags			Admin
//	@Produce		json
//	@Param			page	query		int	false	"Page"	default(1)
//	@Param			limit	query		int	false	"Limit"	default(20)
//	@Success		200		{object}	dto.PaginatedAdminUsersResponse
//	@Failure		401		{object}	map[string]string
//	@Failure		403		{object}	map[string]string
//	@Router			/admin/users [get]
func AdminListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	result, err := services.AdminGetAllUsers(page, limit)
	if err != nil {
		slog.Error("Admin failed to list users", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch users"})
		return
	}

	c.JSON(http.StatusOK, result)
}

// AdminCreateUser godoc
//
//	@Summary		Create a user
//	@Description	Admin — create a new user account
//	@Tags			Admin
//	@Accept			json
//	@Produce		json
//	@Param			user	body		dto.AdminCreateUserRequest	true	"User data"
//	@Success		201		{object}	dto.UserResponse
//	@Failure		409		{object}	map[string]string
//	@Failure		422		{object}	map[string]interface{}
//	@Router			/admin/users [post]
func AdminCreateUser(c *gin.Context) {
	var req dto.AdminCreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := validate.Struct(req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": formatValidationErrors(err)})
		return
	}

	user, err := services.AdminCreateUser(req)
	if err != nil {
		slog.Warn("Admin user creation failed", "error", err.Error())
		switch {
		case errors.Is(err, services.ErrEmailInUse),
			errors.Is(err, services.ErrUsernameInUse):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		}
		return
	}

	slog.Info("Admin created user", "user_id", user.UserID)
	c.JSON(http.StatusCreated, gin.H{"data": user})
}

// AdminUpdateUser godoc
//
//	@Summary		Update a user
//	@Description	Admin — update user details or active status
//	@Tags			Admin
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string						true	"User UUID"
//	@Param			user	body		dto.AdminUpdateUserRequest	true	"Update data"
//	@Success		200		{object}	dto.UserResponse
//	@Failure		404		{object}	map[string]string
//	@Router			/admin/users/{id} [put]
func AdminUpdateUser(c *gin.Context) {
	userID := c.Param("id")

	var req dto.AdminUpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := validate.Struct(req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": formatValidationErrors(err)})
		return
	}

	user, err := services.AdminUpdateUser(userID, req)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrUserNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, services.ErrUsernameInUse):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		case errors.Is(err, services.ErrNoFieldsToUpdate):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user"})
		}
		return
	}

	slog.Info("Admin updated user", "user_id", userID)
	c.JSON(http.StatusOK, gin.H{"data": user})
}

// AdminDeactivateUser godoc
//
//	@Summary		Deactivate a user
//	@Description	Admin — disable a user account
//	@Tags			Admin
//	@Produce		json
//	@Param			id	path		string	true	"User UUID"
//	@Success		200	{object}	map[string]string
//	@Failure		404	{object}	map[string]string
//	@Router			/admin/users/{id}/deactivate [post]
func AdminDeactivateUser(c *gin.Context) {
	userID := c.Param("id")

	if err := services.AdminDeactivateUser(userID); err != nil {
		switch {
		case errors.Is(err, services.ErrUserNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to deactivate user"})
		}
		return
	}

	slog.Info("Admin deactivated user", "user_id", userID)
	c.JSON(http.StatusOK, gin.H{"message": "user deactivated successfully"})
}

// AdminActivateUser godoc
//
//	@Summary		Activate a user
//	@Description	Admin — re-enable a user account
//	@Tags			Admin
//	@Produce		json
//	@Param			id	path		string	true	"User UUID"
//	@Success		200	{object}	map[string]string
//	@Failure		404	{object}	map[string]string
//	@Router			/admin/users/{id}/activate [post]
func AdminActivateUser(c *gin.Context) {
	userID := c.Param("id")

	if err := services.AdminActivateUser(userID); err != nil {
		switch {
		case errors.Is(err, services.ErrUserNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to activate user"})
		}
		return
	}

	slog.Info("Admin activated user", "user_id", userID)
	c.JSON(http.StatusOK, gin.H{"message": "user activated successfully"})
}

// ─────────────────────────────────────────────
// MERCHANTS
// ─────────────────────────────────────────────

// AdminListMerchants godoc
//
//	@Summary		List all merchants
//	@Description	Admin — retrieve all merchants including inactive
//	@Tags			Admin
//	@Produce		json
//	@Param			page	query		int	false	"Page"	default(1)
//	@Param			limit	query		int	false	"Limit"	default(20)
//	@Success		200		{object}	dto.PaginatedAdminMerchantsResponse
//	@Router			/admin/merchants [get]
func AdminListMerchants(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	result, err := services.AdminGetAllMerchants(page, limit)
	if err != nil {
		slog.Error("Admin failed to list merchants", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch merchants"})
		return
	}

	c.JSON(http.StatusOK, result)
}

// AdminCreateMerchant godoc
//
//	@Summary		Create a merchant
//	@Description	Admin — create a new merchant account
//	@Tags			Admin
//	@Accept			json
//	@Produce		json
//	@Param			merchant	body		dto.AdminCreateMerchantRequest	true	"Merchant data"
//	@Success		201			{object}	dto.MerchantResponse
//	@Failure		409			{object}	map[string]string
//	@Failure		422			{object}	map[string]interface{}
//	@Router			/admin/merchants [post]
func AdminCreateMerchant(c *gin.Context) {
	var req dto.AdminCreateMerchantRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := validate.Struct(req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": formatValidationErrors(err)})
		return
	}

	merchant, err := services.AdminCreateMerchant(req)
	if err != nil {
		slog.Warn("Admin merchant creation failed", "error", err.Error())
		switch {
		case errors.Is(err, services.ErrEmailInUse):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create merchant"})
		}
		return
	}

	slog.Info("Admin created merchant", "merchant_id", merchant.MerchantID)
	c.JSON(http.StatusCreated, gin.H{"data": merchant})
}

// AdminUpdateMerchant godoc
//
//	@Summary		Update a merchant
//	@Description	Admin — update merchant details or active status
//	@Tags			Admin
//	@Accept			json
//	@Produce		json
//	@Param			id			path		string							true	"Merchant UUID"
//	@Param			merchant	body		dto.AdminUpdateMerchantRequest	true	"Update data"
//	@Success		200			{object}	dto.MerchantResponse
//	@Failure		404			{object}	map[string]string
//	@Router			/admin/merchants/{id} [put]
func AdminUpdateMerchant(c *gin.Context) {
	merchantID := c.Param("id")

	var req dto.AdminUpdateMerchantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := validate.Struct(req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": formatValidationErrors(err)})
		return
	}

	merchant, err := services.AdminUpdateMerchant(merchantID, req)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrMerchantNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, services.ErrNoFieldsToUpdate):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update merchant"})
		}
		return
	}

	slog.Info("Admin updated merchant", "merchant_id", merchantID)
	c.JSON(http.StatusOK, gin.H{"data": merchant})
}

// AdminDeactivateMerchant godoc
//
//	@Summary		Deactivate a merchant
//	@Description	Admin — disable a merchant account and all its reps
//	@Tags			Admin
//	@Produce		json
//	@Param			id	path		string	true	"Merchant UUID"
//	@Success		200	{object}	map[string]string
//	@Failure		404	{object}	map[string]string
//	@Router			/admin/merchants/{id}/deactivate [post]
func AdminDeactivateMerchant(c *gin.Context) {
	merchantID := c.Param("id")

	if err := services.AdminDeactivateMerchant(merchantID); err != nil {
		switch {
		case errors.Is(err, services.ErrMerchantNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to deactivate merchant"})
		}
		return
	}

	slog.Info("Admin deactivated merchant", "merchant_id", merchantID)
	c.JSON(http.StatusOK, gin.H{"message": "merchant deactivated successfully"})
}

// AdminActivateMerchant godoc
//
//	@Summary		Activate a merchant
//	@Description	Admin — re-enable a merchant account
//	@Tags			Admin
//	@Produce		json
//	@Param			id	path		string	true	"Merchant UUID"
//	@Success		200	{object}	map[string]string
//	@Failure		404	{object}	map[string]string
//	@Router			/admin/merchants/{id}/activate [post]
func AdminActivateMerchant(c *gin.Context) {
	merchantID := c.Param("id")

	if err := services.AdminActivateMerchant(merchantID); err != nil {
		switch {
		case errors.Is(err, services.ErrMerchantNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to activate merchant"})
		}
		return
	}

	slog.Info("Admin activated merchant", "merchant_id", merchantID)
	c.JSON(http.StatusOK, gin.H{"message": "merchant activated successfully"})
}
