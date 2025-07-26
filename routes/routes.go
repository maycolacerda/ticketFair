package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/maycolacerda/ticketfair/controllers"
	_ "github.com/maycolacerda/ticketfair/docs" // Import the generated
	swaggerFiles "github.com/swaggo/files"      // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger"
)

func HandleRequests() {

	r := gin.Default()
	r.GET("/", controllers.GetHome)
	r.NoRoute(controllers.NotFound)

	r.GET("/health", controllers.HealthCheck)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run(":8000") // Listen and serve on localhost:8000
}
