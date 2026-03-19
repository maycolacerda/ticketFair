// database/db.go
package database

import (
	"fmt"
	"log/slog"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() {
	dsn := buildDSN()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:                                   logger.Default.LogMode(logger.Silent),
		DisableNestedTransaction:                 true,
		DisableForeignKeyConstraintWhenMigrating: true,
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

	//migrate(db)

	DB = db
}

func buildDSN() string {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("COCKROACH_USER")
	dbname := os.Getenv("COCKROACH_DB")

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
	if dbname == "" {
		missing = append(missing, "COCKROACH_DB")
	}

	if len(missing) > 0 {
		slog.Error("Missing required database environment variables", "vars", missing)
		os.Exit(1)
	}

	return fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s sslmode=disable",
		host, port, user, dbname,
	)
}

/*func migrate(db *gorm.DB) {
	slog.Info("Running migrations...")

	// Order matters — parent tables must be created before children
	models := []interface{}{
		&models.User{},        // no dependencies
		&models.Merchant{},    // no dependencies
		&models.Profile{},     // depends on User
		&models.Address{},     // depends on Profile
		&models.MerchantRep{}, // depends on Merchant
		&models.Event{},       // depends on Merchant
		&models.Transaction{}, // depends on User and Event
	}

	for _, model := range models {
		if err := db.AutoMigrate(model); err != nil {
			slog.Error("Migration failed", "model", fmt.Sprintf("%T", model), "error", err.Error())
			os.Exit(1)
		}
	}

	slog.Info("Migrations complete")
}
*/
