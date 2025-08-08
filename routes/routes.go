package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/maycolacerda/ticketfair/controllers"
	_ "github.com/maycolacerda/ticketfair/docs" // Import the generated
	"github.com/maycolacerda/ticketfair/middlewares"
	"github.com/maycolacerda/ticketfair/services"
	swaggerFiles "github.com/swaggo/files" // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger"
)

func HandleRequests() {

	r := gin.Default()
	r.GET("/", controllers.GetHome)
	r.NoRoute(controllers.NotFound)

	public := r.Group("/public")
	private := r.Group("/private")

	//public
	public.GET("/health", controllers.HealthCheck)
	public.POST("/register", controllers.NewUser)
	public.POST("/auth/login", services.NewAuthRequest)
	public.POST("/newuser", controllers.NewUser)

	//private
	private.Use(middlewares.AuthMiddleware())
	private.GET("/users", controllers.GetUsers)
	private.GET("/users/:id", controllers.GetUserByID)      //temporário. Retirar depois
	private.GET("/users/me", controllers.CurrentUser)       // Get current user
	private.POST("/profile/new", controllers.CreateProfile) // Create a new profile
	private.POST("/logout", services.Logout)                // Logout user
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run(":8000") // Listen and serve on localhost:8000
}
