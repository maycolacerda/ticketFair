// main.go
package main

import (
	"log/slog"

	"github.com/maycolacerda/ticketfair/database"
	_ "github.com/maycolacerda/ticketfair/docs"
	"github.com/maycolacerda/ticketfair/routes"
	"github.com/maycolacerda/ticketfair/services"
)

//	@title			TicketFair API
//	@version		1.0
//	@description	TicketFair event ticketing platform API

//	@contact.name	TicketFair Support
//	@contact.email	support@ticketfair.com

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		ticketfair.localhost
//	@BasePath	/api/v1

//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				JWT Bearer token. Format: "Bearer <token>"

func main() {
	services.InitLogger()

	slog.Info("Starting TicketFair API")

	slog.Info("Initializing database...")
	database.InitDB()

	slog.Info("Initializing routes...")
	routes.HandleRequests()
}
