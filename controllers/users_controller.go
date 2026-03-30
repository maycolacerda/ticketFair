// controllers/users.go
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

// NewUser godoc
//
//	@Summary		Register a new user
//	@Description	Create a new user with email, password, and username
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			user	body		dto.CreateUserRequest	true	"User registration data"
//	@Success		201		{object}	dto.UserResponse
//	@Failure		400		{object}	map[string]string
//	@Failure		409		{object}	map[string]string
//	@Failure		422		{object}	map[string]interface{}
//	@Router			/public/auth/register [post]
func NewUser(c *gin.Context) {
	var req dto.CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Warn("Invalid request body", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := validate.Struct(req); err != nil {
		errs := formatValidationErrors(err)
		slog.Warn("Validation failed", "errors", errs)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": errs})
		return
	}

	user, err := services.CreateUser(req)
	if err != nil {
		slog.Warn("User creation failed", "error", err.Error())
		switch {
		case errors.Is(err, services.ErrEmailInUse),
			errors.Is(err, services.ErrUsernameInUse):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		}
		return
	}
	slog.Info("User created", "user_id", user.UserID)
	c.JSON(http.StatusCreated, gin.H{"data": user})
}

// GetUsers godoc
//
//	@Summary		List all users
//	@Description	Retrieve a paginated list of users
//	@Tags			Users
//	@Produce		json
//	@Param			page	query		int	false	"Page number"	default(1)
//	@Param			limit	query		int	false	"Page size"		default(20)
//	@Success		200		{object}	dto.PaginatedUsersResponse
//	@Failure		500		{object}	map[string]string
//	@Router			/private/users [get]
func GetUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	result, err := services.GetAllUsers(page, limit)
	if err != nil {
		slog.Error("Failed to fetch users", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch users"})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetUserByID godoc
//
//	@Summary		Get a user by ID
//	@Description	Retrieve a user by their UUID
//	@Tags			Users
//	@Produce		json
//	@Param			id	path		string	true	"User UUID"
//	@Success		200	{object}	dto.UserResponse
//	@Failure		401	{object}	map[string]string
//	@Failure		404	{object}	map[string]string
//	@Router			/private/users/{id} [get]
func GetUserByID(c *gin.Context) {
	userID := c.Param("id")

	user, err := services.GetUserByID(userID)
	if err != nil {
		slog.Warn("User not found", "user_id", userID)
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

// CurrentUser godoc
//
//	@Summary		Get the current authenticated user
//	@Description	Retrieve profile of the currently authenticated user
//	@Tags			Users
//	@Produce		json
//	@Success		200	{object}	dto.UserResponse
//	@Failure		401	{object}	map[string]string
//	@Failure		404	{object}	map[string]string
//	@Router			/private/users/me [get]
func CurrentUser(c *gin.Context) {
	userID, err := services.ExtractTokenID(c)
	if err != nil {
		slog.Warn("Unauthorized token extraction failed")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	user, err := services.GetUserByID(userID)
	if err != nil {
		slog.Warn("User not found", "user_id", userID)
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}
