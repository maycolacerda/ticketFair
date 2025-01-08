package services

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/maycolacerda/ticketfair/models"
)

// todo: fix the username and password check
func LoginHandler(c *gin.Context) {
	var u models.User
	json.NewDecoder(c.Request.Body).Decode(&u)
	fmt.Printf("The user request value %v", u)

	if u.Username == "Chek" && u.Password == "123456" {
		tokenString, err := CreateToken(u.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to generate token",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"token": tokenString,
		})
		return

	} else {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid credentials",
		})
		return

	}
}
