// The above code is a set of controller functions for handling CRUD operations on a Company model in a
// Go web application using the Gin framework.
package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/maycolacerda/ticketfair/database"
	"github.com/maycolacerda/ticketfair/models"
)

func GetCompanies(c *gin.Context) {

	var companies []models.Company
	database.DB.Find(&companies)
	c.JSON(200, companies)

}

func GetCompany(c *gin.Context) {

	var company models.Company

	if err := database.DB.Where("company_id = ?", c.Param("id")).First(&company).Error; err != nil {
		c.JSON(404, gin.H{
			"Message": "Company not found!",
		})
		return
	}

	c.JSON(200, company)
}

func CreateCompany(c *gin.Context) {

	var company models.Company

	if err := c.ShouldBindJSON(&company); err != nil {
		c.JSON(400, gin.H{
			"Message": "Company data is invalid!",
		})
		return
	}

	database.DB.Create(&company)
	c.JSON(201, company)
}

func UpdateCompany(c *gin.Context) {

	var company models.Company

	if err := database.DB.Where("company_id = ?", c.Param("id")).First(&company).Error; err != nil {
		c.JSON(404, gin.H{
			"Message": "Company not found!",
		})
		return
	}

	if err := c.ShouldBindJSON(&company); err != nil {
		c.JSON(400, gin.H{
			"Message": "Company data is invalid!",
		})
		return
	}

	database.DB.Save(&company)
	c.JSON(200, company)
}
