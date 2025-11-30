package main

import (
	"log/slog"

	"github.com/maycolacerda/ticketfair/database"
	"github.com/maycolacerda/ticketfair/routes"
	"github.com/maycolacerda/ticketfair/services"
	_ "github.com/swaggo/files"
	_ "github.com/swaggo/gin-swagger"
)

//	@title			Swagger Example API
//	@version		1.0
//	@description	This is a sample server celler server.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:8000
//	@BasePath	/

//	@securityDefinitions.basic	BasicAuth

// @externalDocs.description	OpenAPI
// @externalDocs.url			https://swagger.io/resources/open-api/
func main() {

	services.Log()
	slog.Info("Initializing database...")
	database.InitDB()
	slog.Info("Initializing routes...")
	routes.HandleRequests()
}
