package controllers

import "github.com/gin-gonic/gin"

func GetUsers(c *gin.Context) {
	c.JSON(401, gin.H{
		"Message": "Users listing not allowed!",
	})

}
func GetUser(c *gin.Context) {
	c.JSON(404, gin.H{
		"Message": "User not found!",
	})

}
