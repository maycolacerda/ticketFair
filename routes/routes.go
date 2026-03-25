// routes/routes.go
package routes

import (
	"io"

	"github.com/gin-gonic/gin"
	"github.com/maycolacerda/ticketfair/controllers"
	"github.com/maycolacerda/ticketfair/middlewares"
	"github.com/maycolacerda/ticketfair/services"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func HandleRequests() {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = io.Discard

	r := gin.Default()

	// Base
	r.GET("/", controllers.GetHome)
	r.NoRoute(controllers.NotFound)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api/v1")

	setupPublicRoutes(api)
	setupPrivateRoutes(api)
	setupMerchantRoutes(api)
	setupAdminRoutes(api)

	r.Run(":8000")
}

func setupPublicRoutes(rg *gin.RouterGroup) {
	public := rg.Group("/public")
	public.Use(middlewares.PublicMiddleware())
	{
		public.GET("/health", controllers.HealthCheck)

		auth := public.Group("/auth")
		{
			auth.POST("/register", controllers.NewUser)
			auth.POST("/client/login", controllers.ClientLogin)
			auth.POST("/merchant/login", controllers.MerchantLogin)
			auth.POST("/rep/login", controllers.MerchantRepLogin)
			auth.POST("/logout", controllers.Logout)
		}

		merchant := public.Group("/merchant")
		{
			merchant.POST("/register", controllers.NewMerchant)
		}

		events := public.Group("/events")
		{
			events.GET("/", controllers.GetEvents)
			events.GET("/:id", controllers.GetEventByID)
		}
	}
}

func setupPrivateRoutes(rg *gin.RouterGroup) {
	private := rg.Group("/private")
	private.Use(middlewares.ClientMiddleware())
	{
		users := private.Group("/users")
		{
			users.GET("/", controllers.GetUsers)
			users.GET("/me", controllers.CurrentUser)
			users.GET("/:id", controllers.GetUserByID)
		}

		profile := private.Group("/profile")
		{
			profile.GET("/", controllers.GetProfile)
			profile.POST("/", controllers.CreateProfile)
			profile.PUT("/", controllers.UpdateProfile)
		}

		tickets := private.Group("/tickets")
		{
			tickets.GET("/", controllers.GetMyTickets)
			tickets.GET("/:id", controllers.GetTicketByID)
			tickets.POST("/purchase", controllers.PurchaseTicket)
			tickets.POST("/refund", controllers.RefundTicket)
		}

		transactions := private.Group("/transactions")
		{
			transactions.GET("/", controllers.GetMyTransactions)
		}
		verify := private.Group("/verify")
		{
			verify.POST("/email/send", controllers.SendEmailVerification)
			verify.POST("/email", controllers.VerifyEmail)
			verify.POST("/phone/send", controllers.SendPhoneVerification)
			verify.POST("/phone", controllers.VerifyPhone)
		}

		private.POST("/logout", controllers.Logout)
	}
}
func setupMerchantRoutes(rg *gin.RouterGroup) {
	merchant := rg.Group("/merchant")
	merchant.Use(middlewares.MerchantMiddleware())
	{
		merchant.PUT("/update", controllers.UpdateMerchant)
		merchant.POST("/logout", controllers.Logout)

		events := merchant.Group("/events")
		{
			events.POST("/new", controllers.NewEvent)
			events.PUT("/:id", controllers.UpdateEvent)
		}

		rep := merchant.Group("/rep")
		rep.Use(middlewares.MerchantRepMiddleware(services.RoleMerchantAdmin))
		{
			rep.POST("/new", controllers.NewMerchantRep)
			rep.PUT("/:id", controllers.UpdateMerchantRep)
		}
	}
}

func setupAdminRoutes(rg *gin.RouterGroup) {

	adminPublic := rg.Group("/admin")
	{
		adminPublic.POST("/auth/login", controllers.AdminLogin)
	}

	admin := rg.Group("/admin")
	admin.Use(middlewares.AdminMiddleware())
	{
		users := admin.Group("/users")
		{
			users.GET("/", controllers.AdminListUsers)
			users.POST("/", controllers.AdminCreateUser)
			users.PUT("/:id", controllers.AdminUpdateUser)
			users.POST("/:id/deactivate", controllers.AdminDeactivateUser)
			users.POST("/:id/activate", controllers.AdminActivateUser)
		}

		merchants := admin.Group("/merchants")
		{
			merchants.GET("/", controllers.AdminListMerchants)
			merchants.POST("/", controllers.AdminCreateMerchant)
			merchants.PUT("/:id", controllers.AdminUpdateMerchant)
			merchants.POST("/:id/deactivate", controllers.AdminDeactivateMerchant)
			merchants.POST("/:id/activate", controllers.AdminActivateMerchant)
		}
	}
}
