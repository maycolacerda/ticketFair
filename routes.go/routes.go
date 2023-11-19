package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/maycolacerda/ticketfair/controllers"
)

func Handler() {

	r := gin.Default()

	r.GET("/", controllers.HelloWorld)
	r.GET("/index", controllers.HelloWorld)
	r.GET("/tickets", controllers.GetTickets)
	r.GET("/users", controllers.GetUsers)
	r.GET("/user/:user_id", controllers.GetUser)

	http.ListenAndServe("localhost:8080", r)
}
