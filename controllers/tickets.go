package controllers

import "github.com/gin-gonic/gin"

func GetTickets(c *gin.Context) {
	c.JSON(404, gin.H{
		"Message": "Tickets not found yet!",
	})

}
