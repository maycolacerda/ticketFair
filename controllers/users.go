package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/maycolacerda/ticketfair/database"
	"github.com/maycolacerda/ticketfair/models"
)

// NewUser godoc
//
//	@Summary		Create a new user.
//	@Description	Create a new user with email, password, and username.
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			user	body	models.User	true	"User data"
//	@Success		200	{object}	map[string]string
//	@Failure		400	{object}	map[string]string
//	@Router			/users/new [post]
func NewUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err := user.Validate()
	if len(err) > 0 {
		c.JSON(400, gin.H{"error": err})
		return
	} else {
		database.DB.Create(&user)
		c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})

	}
}

// GetUsers godoc
//
//	@Summary		Get all users.
//	@Description	Retrieve a list of all users.
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		models.User
//	@Router			/users [get]
func GetUsers(c *gin.Context) {
	var users []models.User
	database.DB.Find(&users)
	c.JSON(http.StatusOK, users)
}

// GetUserByID godoc
//
//	@Summary		Get a user by ID.
//	@Description	Retrieve a user by their ID.
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"User ID"
//	@Success		200	{object}	models.User
//	@Failure		404	{object}	map[string]string
//	@Router			/users/{id} [get]
func GetUserByID(c *gin.Context) {
	userID := c.Param("id")
	var user models.User
	if err := database.DB.First(&user, "user_id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}
