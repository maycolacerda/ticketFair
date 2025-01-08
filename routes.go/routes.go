package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/maycolacerda/ticketfair/controllers"
	"github.com/maycolacerda/ticketfair/services"
)

func HandlerFunc() {

	r := gin.Default()

	protected := r.Group("/protected")
	{

		protected.GET("/tickets", controllers.GetTickets)
		protected.GET("/users", controllers.GetUsers)
		protected.GET("/user/:user_id", controllers.GetUser)
	}
	open := r.Group("/open")
	{
		open.GET("/", controllers.HelloWorld)
		open.GET("/index", controllers.HelloWorld)
		open.POST("/login", services.LoginHandler)
	}

	http.ListenAndServe("localhost:8080", r)
}
