// tests/integration/auth_test.go
package integration

import (
	"net/http"
	"testing"

	"github.com/maycolacerda/ticketfair/services"
	"github.com/maycolacerda/ticketfair/tests/testutil"
)

func TestAuth_Register(t *testing.T) {
	db := testutil.SetupTestDB(t)
	router := testutil.NewRouter()
	t.Cleanup(func() { testutil.TruncateAll(db) })

	t.Run("success", func(t *testing.T) {
		w := testutil.POST(router, "/api/v1/public/auth/register").
			WithBody(testutil.J(
				"email", "newuser@example.com",
				"username", "newuser",
				"password", "PassW0rd!",
			)).Do()

		if w.Code != http.StatusCreated {
			t.Fatalf("expected 201, got %d — body: %s", w.Code, w.Body.String())
		}

		data := testutil.DataField(w)
		if data["user_id"] == "" {
			t.Error("expected user_id in response")
		}
		if data["password"] != nil {
			t.Error("password must not be exposed in response")
		}
	})

	t.Run("duplicate email returns 409", func(t *testing.T) {
		testutil.CreateUser(db, "dup@example.com", "dupuser")

		w := testutil.POST(router, "/api/v1/public/auth/register").
			WithBody(testutil.J(
				"email", "dup@example.com",
				"username", "dupuser2",
				"password", "PassW0rd!",
			)).Do()

		if w.Code != http.StatusConflict {
			t.Fatalf("expected 409, got %d", w.Code)
		}
	})

	t.Run("duplicate username returns 409", func(t *testing.T) {
		testutil.CreateUser(db, "another@example.com", "takenuser")

		w := testutil.POST(router, "/api/v1/public/auth/register").
			WithBody(testutil.J(
				"email", "fresh@example.com",
				"username", "takenuser",
				"password", "PassW0rd!",
			)).Do()

		if w.Code != http.StatusConflict {
			t.Fatalf("expected 409, got %d", w.Code)
		}
	})

	t.Run("weak password returns 422", func(t *testing.T) {
		w := testutil.POST(router, "/api/v1/public/auth/register").
			WithBody(testutil.J(
				"email", "weak@example.com",
				"username", "weakuser",
				"password", "weak",
			)).Do()

		if w.Code != http.StatusUnprocessableEntity {
			t.Fatalf("expected 422, got %d", w.Code)
		}
	})

	t.Run("invalid email returns 422", func(t *testing.T) {
		w := testutil.POST(router, "/api/v1/public/auth/register").
			WithBody(testutil.J(
				"email", "notanemail",
				"username", "someuser",
				"password", "PassW0rd!",
			)).Do()

		if w.Code != http.StatusUnprocessableEntity {
			t.Fatalf("expected 422, got %d", w.Code)
		}
	})
}

func TestAuth_ClientLogin(t *testing.T) {
	db := testutil.SetupTestDB(t)
	router := testutil.NewRouter()
	t.Cleanup(func() { testutil.TruncateAll(db) })

	testutil.CreateUser(db, "login@example.com", "loginuser")

	t.Run("success returns token", func(t *testing.T) {
		w := testutil.POST(router, "/api/v1/public/auth/client/login").
			WithBody(testutil.J(
				"email", "login@example.com",
				"password", testutil.TestPassword,
			)).Do()

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d — body: %s", w.Code, w.Body.String())
		}

		data := testutil.DataField(w)
		if data == nil || data["token"] == "" {
			t.Fatal("expected token in response")
		}
		if data["user"] == nil {
			t.Fatal("expected user in response")
		}
	})

	t.Run("wrong password returns 401", func(t *testing.T) {
		w := testutil.POST(router, "/api/v1/public/auth/client/login").
			WithBody(testutil.J(
				"email", "login@example.com",
				"password", "WrongPass!",
			)).Do()

		if w.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401, got %d", w.Code)
		}
	})

	t.Run("unknown email returns 401", func(t *testing.T) {
		w := testutil.POST(router, "/api/v1/public/auth/client/login").
			WithBody(testutil.J(
				"email", "ghost@example.com",
				"password", testutil.TestPassword,
			)).Do()

		if w.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401, got %d", w.Code)
		}
	})

	t.Run("disabled account returns 403", func(t *testing.T) {
		u := testutil.CreateUser(db, "disabled@example.com", "disableduser")
		db.Model(u).Update("active", false)

		w := testutil.POST(router, "/api/v1/public/auth/client/login").
			WithBody(testutil.J(
				"email", "disabled@example.com",
				"password", testutil.TestPassword,
			)).Do()

		if w.Code != http.StatusForbidden {
			t.Fatalf("expected 403, got %d", w.Code)
		}
	})
}

func TestAuth_MerchantLogin(t *testing.T) {
	db := testutil.SetupTestDB(t)
	router := testutil.NewRouter()
	t.Cleanup(func() { testutil.TruncateAll(db) })

	testutil.CreateMerchant(db, "merchant@example.com")

	t.Run("success", func(t *testing.T) {
		w := testutil.POST(router, "/api/v1/public/auth/merchant/login").
			WithBody(testutil.J(
				"email", "merchant@example.com",
				"password", testutil.TestPassword,
			)).Do()

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d — body: %s", w.Code, w.Body.String())
		}
		if testutil.DataString(w, "token") == "" {
			t.Fatal("expected token")
		}
	})

	t.Run("wrong password returns 401", func(t *testing.T) {
		w := testutil.POST(router, "/api/v1/public/auth/merchant/login").
			WithBody(testutil.J(
				"email", "merchant@example.com",
				"password", "wrong",
			)).Do()

		if w.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401, got %d", w.Code)
		}
	})
}

func TestAuth_PasswordReset(t *testing.T) {
	db := testutil.SetupTestDB(t)
	router := testutil.NewRouter()
	t.Cleanup(func() { testutil.TruncateAll(db) })

	testutil.CreateUser(db, "reset@example.com", "resetuser")

	t.Run("forgot password always returns 200", func(t *testing.T) {
		// Known email
		w := testutil.POST(router, "/api/v1/public/auth/password/forgot").
			WithBody(testutil.J("email", "reset@example.com")).
			Do()

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}

		// Unknown email — same response (prevent enumeration)
		w2 := testutil.POST(router, "/api/v1/public/auth/password/forgot").
			WithBody(testutil.J("email", "ghost@example.com")).
			Do()

		if w2.Code != http.StatusOK {
			t.Fatalf("expected 200 for unknown email, got %d", w2.Code)
		}
	})

	t.Run("reset with wrong code returns 400", func(t *testing.T) {
		// Trigger forgot first to create a code
		testutil.POST(router, "/api/v1/public/auth/password/forgot").
			WithBody(testutil.J("email", "reset@example.com")).
			Do()

		w := testutil.POST(router, "/api/v1/public/auth/password/reset").
			WithBody(testutil.J(
				"email", "reset@example.com",
				"code", "000000",
				"new_password", "NewPassW0rd!",
			)).Do()

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d — body: %s", w.Code, w.Body.String())
		}
	})

	t.Run("reset with correct code succeeds and invalidates old password", func(t *testing.T) {
		// Get the reset code directly from DB
		testutil.POST(router, "/api/v1/public/auth/password/forgot").
			WithBody(testutil.J("email", "reset@example.com")).
			Do()

		var code string
		db.Raw("SELECT code FROM password_resets WHERE used_at IS NULL AND deleted_at IS NULL ORDER BY created_at DESC LIMIT 1").Scan(&code)

		if code == "" {
			t.Fatal("expected a reset code in DB")
		}

		// Reset the password
		w := testutil.POST(router, "/api/v1/public/auth/password/reset").
			WithBody(testutil.J(
				"email", "reset@example.com",
				"code", code,
				"new_password", "BrandNewP@ss1",
			)).Do()

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d — body: %s", w.Code, w.Body.String())
		}

		// Old password should fail
		w2 := testutil.POST(router, "/api/v1/public/auth/client/login").
			WithBody(testutil.J(
				"email", "reset@example.com",
				"password", testutil.TestPassword,
			)).Do()
		if w2.Code != http.StatusUnauthorized {
			t.Error("old password should no longer work")
		}

		// New password should succeed
		w3 := testutil.POST(router, "/api/v1/public/auth/client/login").
			WithBody(testutil.J(
				"email", "reset@example.com",
				"password", "BrandNewP@ss1",
			)).Do()
		if w3.Code != http.StatusOK {
			t.Errorf("new password should work, got %d", w3.Code)
		}
	})

	t.Run("code cannot be reused", func(t *testing.T) {
		testutil.POST(router, "/api/v1/public/auth/password/forgot").
			WithBody(testutil.J("email", "reset@example.com")).
			Do()

		var code string
		db.Raw("SELECT code FROM password_resets WHERE used_at IS NULL AND deleted_at IS NULL ORDER BY created_at DESC LIMIT 1").Scan(&code)

		body := testutil.J(
			"email", "reset@example.com",
			"code", code,
			"new_password", "AnotherP@ss1",
		)

		// First use — OK
		testutil.POST(router, "/api/v1/public/auth/password/reset").WithBody(body).Do()

		// Second use — should fail
		w := testutil.POST(router, "/api/v1/public/auth/password/reset").WithBody(body).Do()
		if w.Code != http.StatusBadRequest && w.Code != http.StatusNotFound {
			t.Errorf("expected 400 or 404 on reuse, got %d", w.Code)
		}
	})
}

func TestAuth_RepLogin(t *testing.T) {
	db := testutil.SetupTestDB(t)
	router := testutil.NewRouter()
	t.Cleanup(func() { testutil.TruncateAll(db) })

	m := testutil.CreateMerchant(db, "repmerchant@example.com")
	testutil.CreateMerchantRep(db, m.MerchantID, "rep@example.com", "admin")

	t.Run("rep login success", func(t *testing.T) {
		w := testutil.POST(router, "/api/v1/public/auth/rep/login").
			WithBody(testutil.J(
				"email", "rep@example.com",
				"password", testutil.TestPassword,
			)).Do()

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d — body: %s", w.Code, w.Body.String())
		}
	})

	t.Run("rep login blocked when merchant disabled", func(t *testing.T) {
		db.Model(m).Update("active", false)

		w := testutil.POST(router, "/api/v1/public/auth/rep/login").
			WithBody(testutil.J(
				"email", "rep@example.com",
				"password", testutil.TestPassword,
			)).Do()

		if w.Code != http.StatusForbidden {
			t.Fatalf("expected 403, got %d", w.Code)
		}
		// Re-enable for cleanup
		db.Model(m).Update("active", true)
	})
}

// loginAs is a helper used across integration tests to get a token.
func loginAs(router interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}, email, password string) string {
	w := testutil.POST(testutil.RouterFromEngine(router), "/api/v1/public/auth/client/login").
		WithBody(testutil.J("email", email, "password", password)).
		Do()
	return testutil.DataString(w, "token")
}

// loginMerchant returns a merchant token.
func loginMerchant(router interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}, email string) string {
	var token string
	w := testutil.POST(testutil.RouterFromEngine(router), "/api/v1/public/auth/merchant/login").
		WithBody(testutil.J("email", email, "password", testutil.TestPassword)).
		Do()
	body := testutil.BodyString(w)
	if data, ok := body["data"].(map[string]interface{}); ok {
		token, _ = data["token"].(string)
	}
	return token
}

var _ = services.GenerateToken // ensure services is imported
