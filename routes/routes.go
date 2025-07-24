package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/maycolacerda/ticketfair/controllers"
)

func HandleRequests() {

	r := gin.Default()
	r.GET("/", controllers.GetHome)
	r.NoRoute(controllers.NotFound)

	r.GET("/health", controllers.HealthCheck)

	r.Run(":8000") // Listen and serve on localhost:8000
}
