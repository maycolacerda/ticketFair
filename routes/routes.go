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
	merchant := r.Group("/merchant")

	//public
	public.GET("/health", controllers.HealthCheck)
	public.POST("/register", controllers.NewUser)
	public.POST("/auth/client/login", services.NewAuthRequestClient)
	public.POST("/auth/merchant/login", services.NewAuthRequestMerchant)
	public.POST("/auth/logout", services.Logout)
	public.POST("/newuser", controllers.NewUser)
	public.POST("/merchants/new/merchant", controllers.NewMerchant)
	public.POST("/merchant/new/rep", controllers.NewMerchantRep)

	//private
	private.Use(middlewares.ClientMiddleware())
	private.GET("/users", controllers.GetUsers)
	private.GET("/users/:id", controllers.GetUserByID) //temporário. Retirar depois
	private.GET("/users/me", controllers.CurrentUser)
	private.POST("/profile/new", controllers.CreateProfile)
	private.POST("/profile/update", controllers.UpdateProfile)
	private.POST("/logout", services.Logout)

	//merchant
	merchant.Use(middlewares.MerchantMiddleware())
	merchant.POST("/merchant/events/new", controllers.NewEvent)
	merchant.POST("/merchant/login", services.NewAuthRequestMerchant)
	merchant.POST("/merchant/logout", services.Logout)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run(":8000") // Listen and serve on localhost:8000
}
