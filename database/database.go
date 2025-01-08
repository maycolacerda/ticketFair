package database

import (
	"database/sql"
	"log"

	"github.com/maycolacerda/ticketfair/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB  *gorm.DB
	err error
)

// The Connect function establishes a connection to a MySQL database and performs automatic migration.
func Connect() {
	db, erro := sql.Open("mysql", "root:root@tcp(localhost:5432)/ticketfair?parseTime=true")
	if erro != nil {
		log.Fatal(erro.Error())
	}
	DB, err = gorm.Open(mysql.New(mysql.Config{
		Conn: db,
	}), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})

	if err != nil {
		log.Fatal("erro na conexão com o banco de dados: " + err.Error())
	} else {
		log.Println("conexão com o banco de dados estabelecida com sucesso")
	}
	Automigrate()
}

// The `Automigrate` function is used to automatically migrate database tables for the `User`,
// `Ticket`, `Event`, `TicketGroup`, and `UserProfile` models.
func Automigrate() {
	DB.AutoMigrate(&models.User{}, &models.Ticket{}, &models.Event{}, &models.TicketGroup{}, &models.UserProfile{})
}
