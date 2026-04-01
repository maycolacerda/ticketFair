// tests/integration/main_test.go
package integration

import (
	"os"
	"testing"

	"github.com/maycolacerda/ticketfair/services"
	"github.com/maycolacerda/ticketfair/tests/testutil"
)

func TestMain(m *testing.M) {
	// Wire the token generator into the testutil package so fixtures can create tokens
	testutil.SetTokenGenerator(services.GenerateToken)

	// Set JWT_SECRET for tests
	if os.Getenv("JWT_SECRET") == "" {
		os.Setenv("JWT_SECRET", "test_secret_key_for_integration_tests_only")
	}

	// Set test DB env if not set (CI will override these)
	setDefault("TEST_DB_HOST", "localhost")
	setDefault("TEST_DB_PORT", "26258")
	setDefault("TEST_DB_USER", "root")
	setDefault("TEST_DB_NAME", "ticketfair_test")

	os.Exit(m.Run())
}

func setDefault(key, value string) {
	if os.Getenv(key) == "" {
		os.Setenv(key, value)
	}
}

// services_GenerateToken is exposed for use in test helpers in other files.
var services_GenerateToken = services.GenerateToken
