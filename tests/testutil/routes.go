// tests/testutil/router.go
package testutil

import (
	"io"

	"github.com/gin-gonic/gin"
	"github.com/maycolacerda/ticketfair/controllers"
	"github.com/maycolacerda/ticketfair/middlewares"
	"github.com/maycolacerda/ticketfair/services"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// NewRouter builds the full application router for integration testing.
// Uses the same route registration as production.
func NewRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	r.GET("/", controllers.GetHome)
	r.NoRoute(controllers.NotFound)

	api := r.Group("/api/v1")

	setupPublic(api)
	setupPrivate(api)
	setupMerchant(api)
	setupAdmin(api)

	return r
}

func setupPublic(rg *gin.RouterGroup) {
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
			auth.POST("/password/forgot", controllers.ForgotPassword)
			auth.POST("/password/reset", controllers.ResetPassword)
		}

		public.Group("/merchant").POST("/register", controllers.NewMerchant)

		events := public.Group("/events")
		{
			events.GET("/", controllers.GetEvents)
			events.GET("/:id", controllers.GetEventByID)
			events.GET("/:id/ticket-types", controllers.ListTicketTypes)
		}

		public.POST("/webhooks/stripe", controllers.StripeWebhook)
	}
}

func setupPrivate(rg *gin.RouterGroup) {
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

		verify := private.Group("/verify")
		{
			verify.POST("/email/send", controllers.SendEmailVerification)
			verify.POST("/email", controllers.VerifyEmail)
			verify.POST("/phone/send", controllers.SendPhoneVerification)
			verify.POST("/phone", controllers.VerifyPhone)
		}

		tickets := private.Group("/tickets")
		{
			tickets.GET("/", controllers.GetMyTickets)
			tickets.GET("/:id", controllers.GetTicketByID)
			tickets.POST("/purchase", controllers.PurchaseTicket)
			tickets.POST("/refund", controllers.RefundTicket)
		}

		private.Group("/transactions").GET("/", controllers.GetMyTransactions)

		payments := private.Group("/payments")
		{
			payments.GET("/", controllers.GetMyPayments)
			payments.POST("/intent", controllers.CreatePaymentIntent)
			payments.POST("/:id/refund", controllers.RefundPayment)
		}

		private.POST("/logout", controllers.Logout)
	}
}

func setupMerchant(rg *gin.RouterGroup) {
	merchant := rg.Group("/merchant")
	merchant.Use(middlewares.MerchantMiddleware())
	{
		merchant.PUT("/update", controllers.UpdateMerchant)
		merchant.POST("/logout", controllers.Logout)

		events := merchant.Group("/events")
		{
			events.POST("/new", controllers.NewEvent)
			events.PUT("/:id", controllers.UpdateEvent)
			events.POST("/:id/image", controllers.UploadEventImage)
			events.DELETE("/:id/image", controllers.DeleteEventImageHandler)
			events.GET("/:id/ticket-types", controllers.ListAllTicketTypes)
			events.POST("/:id/ticket-types", controllers.CreateTicketType)
			events.PUT("/:id/ticket-types/:ttid", controllers.UpdateTicketType)
			events.DELETE("/:id/ticket-types/:ttid", controllers.DeleteTicketType)
		}

		tickets := merchant.Group("/tickets")
		tickets.POST("/:id/validate", controllers.ValidateTicket)

		rep := merchant.Group("/rep")
		rep.Use(middlewares.MerchantRepMiddleware(services.RoleMerchantAdmin))
		{
			rep.POST("/new", controllers.NewMerchantRep)
			rep.PUT("/:id", controllers.UpdateMerchantRep)
		}
	}
}

func setupAdmin(rg *gin.RouterGroup) {
	rg.Group("/admin").POST("/auth/login", controllers.AdminLogin)

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

// helper to discard all output in test router
var _ = io.Discard
