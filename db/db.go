package db

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/maycolacerda/ticketfair/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB  *gorm.DB
	err error
)

func InitDB() {

	dsn := "user:password@tcp(localhost:3306)/ticketfair_db?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.New(mysql.Config{
		DSN: dsn,
	}), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}
	log.Println("Database connection established successfully")

	Automigrate()
}

func Automigrate() {
	DB.AutoMigrate(&models.User{})
}
