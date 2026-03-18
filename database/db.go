// database/db.go
package database

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/maycolacerda/ticketfair/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() {
	dsn := buildDSN()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
		// CockroachDB uses implicit transactions internally —
		// this prevents GORM from wrapping every query in a savepoint
		DisableNestedTransaction: true,
	})
	if err != nil {
		slog.Error("Failed to connect to database", "error", err.Error())
		os.Exit(1)
	}

	sqlDB, err := db.DB()
	if err != nil {
		slog.Error("Failed to get sql.DB instance", "error", err.Error())
		os.Exit(1)
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)

	if err := sqlDB.Ping(); err != nil {
		slog.Error("Database ping failed", "error", err.Error())
		os.Exit(1)
	}

	slog.Info("Database connected successfully")

	migrate(db)

	DB = db
}

func buildDSN() string {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("COCKROACH_USER")
	password := os.Getenv("COCKROACH_PASSWORD")
	dbname := os.Getenv("COCKROACH_DB")
	sslmode := os.Getenv("DB_SSLMODE")

	if sslmode == "" {
		sslmode = "disable" // ← for local Docker; use "require" in production
	}

	missing := []string{}
	if host == "" {
		missing = append(missing, "DB_HOST")
	}
	if port == "" {
		missing = append(missing, "DB_PORT")
	}
	if user == "" {
		missing = append(missing, "COCKROACH_USER")
	}
	if password == "" {
		missing = append(missing, "COCKROACH_PASSWORD")
	}
	if dbname == "" {
		missing = append(missing, "COCKROACH_DB")
	}

	if len(missing) > 0 {
		slog.Error("Missing required database environment variables", "vars", missing)
		os.Exit(1)
	}

	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode,
	)
}

func migrate(db *gorm.DB) {
	slog.Info("Running migrations...")

	err := db.AutoMigrate(
		&models.User{},
		&models.Merchant{},
		&models.MerchantRep{},
		&models.Profile{},
		&models.Address{},
		&models.Event{},
		&models.Transaction{},
	)
	if err != nil {
		slog.Error("Migration failed", "error", err.Error())
		os.Exit(1)
	}

	slog.Info("Migrations complete")
}
