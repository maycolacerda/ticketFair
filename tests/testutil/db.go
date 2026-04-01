// tests/testutil/db.go
package testutil

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/maycolacerda/ticketfair/database"
	"github.com/maycolacerda/ticketfair/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	once   sync.Once
	testDB *gorm.DB
)

// SetupTestDB connects to the integration test CockroachDB instance.
// Safe to call from multiple test files — connection is shared via sync.Once.
func SetupTestDB(t interface {
	Helper()
	Fatalf(string, ...any)
}) *gorm.DB {
	t.Helper()

	once.Do(func() {
		db, err := gorm.Open(postgres.Open(buildTestDSN()), &gorm.Config{
			Logger:                                   logger.Default.LogMode(logger.Silent),
			DisableNestedTransaction:                 true,
			DisableForeignKeyConstraintWhenMigrating: true,
		})
		if err != nil {
			log.Fatalf("testutil: failed to connect to test DB: %v", err)
		}

		if err := db.AutoMigrate(
			&models.User{},
			&models.Merchant{},
			&models.MerchantRep{},
			&models.Profile{},
			&models.Address{},
			&models.Event{},
			&models.TicketType{},
			&models.Transaction{},
			&models.Ticket{},
			&models.Verification{},
			&models.PasswordReset{},
			&models.Payment{},
			&models.Admin{},
		); err != nil {
			log.Fatalf("testutil: failed to migrate test DB: %v", err)
		}

		testDB = db
		// Override global DB so all services use the test database
		database.DB = db
	})

	return testDB
}

// TruncateAll deletes all rows from every table.
// Call at the start of each test or in t.Cleanup() for isolation.
func TruncateAll(db *gorm.DB) {
	tables := []string{
		"payments",
		"tickets",
		"transactions",
		"ticket_types",
		"verifications",
		"password_resets",
		"addresses",
		"profiles",
		"events",
		"merchant_reps",
		"merchants",
		"admins",
		"users",
	}
	for _, t := range tables {
		db.Exec(fmt.Sprintf("DELETE FROM %s", t))
	}
}

func buildTestDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s sslmode=disable",
		getEnv("TEST_DB_HOST", "localhost"),
		getEnv("TEST_DB_PORT", "26258"),
		getEnv("TEST_DB_USER", "root"),
		getEnv("TEST_DB_NAME", "ticketfair_test"),
	)
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
