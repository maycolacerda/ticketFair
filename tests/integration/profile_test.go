// tests/integration/profile_test.go
package integration

import (
	"net/http"
	"testing"

	"github.com/maycolacerda/ticketfair/tests/testutil"
)

func TestProfile_CRUD(t *testing.T) {
	db := testutil.SetupTestDB(t)
	router := testutil.NewRouter()
	t.Cleanup(func() { testutil.TruncateAll(db) })

	u := testutil.CreateUser(db, "profile@example.com", "profileuser")
	token, _, _ := services_token(u.UserID, "client", "")

	t.Run("create profile", func(t *testing.T) {
		w := testutil.POST(router, "/api/v1/private/profile/").
			WithAuth(token).
			WithBody(testutil.J(
				"first_name", "João",
				"last_name", "Silva",
				"phone_number", "44999111111",
				"address", map[string]interface{}{
					"street":   "Rua das Flores, 123",
					"city":     "Cianorte",
					"state":    "PR",
					"country":  "BR",
					"zip_code": "87200000",
				},
			)).Do()

		if w.Code != http.StatusCreated {
			t.Fatalf("expected 201, got %d — body: %s", w.Code, w.Body.String())
		}

		data := testutil.DataField(w)
		if data["profile_id"] == "" {
			t.Error("expected profile_id")
		}
		if data["verified_email"] != false {
			t.Error("verified_email should be false initially")
		}
		if data["verified_phone"] != false {
			t.Error("verified_phone should be false initially")
		}
		if addr, ok := data["address"].(map[string]interface{}); ok {
			if addr["city"] != "Cianorte" {
				t.Errorf("expected city Cianorte, got %v", addr["city"])
			}
		} else {
			t.Error("expected address in response")
		}
	})

	t.Run("cannot create duplicate profile", func(t *testing.T) {
		w := testutil.POST(router, "/api/v1/private/profile/").
			WithAuth(token).
			WithBody(testutil.J(
				"first_name", "João",
				"last_name", "Silva",
				"phone_number", "44999222222",
				"address", map[string]interface{}{
					"street":   "Rua B, 456",
					"city":     "Maringá",
					"state":    "PR",
					"country":  "BR",
					"zip_code": "87050000",
				},
			)).Do()

		if w.Code != http.StatusConflict {
			t.Fatalf("expected 409, got %d", w.Code)
		}
	})

	t.Run("get profile", func(t *testing.T) {
		w := testutil.GET(router, "/api/v1/private/profile/").
			WithAuth(token).
			Do()

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}

		data := testutil.DataField(w)
		if data["first_name"] != "João" {
			t.Errorf("expected first_name João, got %v", data["first_name"])
		}
	})

	t.Run("update profile", func(t *testing.T) {
		w := testutil.PUT(router, "/api/v1/private/profile/").
			WithAuth(token).
			WithBody(testutil.J(
				"first_name", "João Updated",
				"address", map[string]interface{}{
					"city": "Maringá",
				},
			)).Do()

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d — body: %s", w.Code, w.Body.String())
		}

		data := testutil.DataField(w)
		if data["first_name"] != "João Updated" {
			t.Errorf("expected updated first_name, got %v", data["first_name"])
		}
		if addr, ok := data["address"].(map[string]interface{}); ok {
			if addr["city"] != "Maringá" {
				t.Errorf("expected city Maringá, got %v", addr["city"])
			}
		}
	})

	t.Run("update with no fields returns 400", func(t *testing.T) {
		w := testutil.PUT(router, "/api/v1/private/profile/").
			WithAuth(token).
			WithBody(testutil.J()).
			Do()

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", w.Code)
		}
	})

	t.Run("phone already in use returns 409", func(t *testing.T) {
		u2 := testutil.CreateUser(db, "profile2@example.com", "profileuser2")
		testutil.CreateProfile(db, u2.UserID, "44999333333")

		w := testutil.PUT(router, "/api/v1/private/profile/").
			WithAuth(token).
			WithBody(testutil.J("phone_number", "44999333333")).
			Do()

		if w.Code != http.StatusConflict {
			t.Fatalf("expected 409, got %d", w.Code)
		}
	})
}

func TestVerification_Email(t *testing.T) {
	db := testutil.SetupTestDB(t)
	router := testutil.NewRouter()
	t.Cleanup(func() { testutil.TruncateAll(db) })

	u := testutil.CreateUser(db, "verify@example.com", "verifyuser")
	testutil.CreateProfile(db, u.UserID, "44999444444")
	token, _, _ := services_token(u.UserID, "client", "")

	t.Run("send email verification", func(t *testing.T) {
		w := testutil.POST(router, "/api/v1/private/verify/email/send").
			WithAuth(token).
			Do()

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d — body: %s", w.Code, w.Body.String())
		}

		data := testutil.DataField(w)
		if data["verified"] != false {
			t.Error("expected verified = false")
		}
		msg, _ := data["message"].(string)
		if msg == "" {
			t.Error("expected message in response")
		}
		// Email must be masked
		if len(msg) > 0 && !containsStr(msg, "***") {
			t.Error("expected masked email in message")
		}
	})

	t.Run("verify with wrong code returns 400", func(t *testing.T) {
		w := testutil.POST(router, "/api/v1/private/verify/email").
			WithAuth(token).
			WithBody(testutil.J("code", "000000")).
			Do()

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", w.Code)
		}
	})

	t.Run("verify with correct code marks email verified", func(t *testing.T) {
		// Get code from DB
		var code string
		db.Raw("SELECT code FROM verifications WHERE user_id = ? AND type = 'email' AND used_at IS NULL AND deleted_at IS NULL ORDER BY created_at DESC LIMIT 1", u.UserID).Scan(&code)

		if code == "" {
			t.Fatal("expected verification code in DB — did send succeed?")
		}

		w := testutil.POST(router, "/api/v1/private/verify/email").
			WithAuth(token).
			WithBody(testutil.J("code", code)).
			Do()

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d — body: %s", w.Code, w.Body.String())
		}

		data := testutil.DataField(w)
		if data["verified"] != true {
			t.Error("expected verified = true")
		}

		// Profile should now show verified_email = true
		pw := testutil.GET(router, "/api/v1/private/profile/").WithAuth(token).Do()
		pdata := testutil.DataField(pw)
		if pdata["verified_email"] != true {
			t.Error("profile.verified_email should be true")
		}
	})

	t.Run("send again after verified returns 409", func(t *testing.T) {
		w := testutil.POST(router, "/api/v1/private/verify/email/send").
			WithAuth(token).
			Do()

		if w.Code != http.StatusConflict {
			t.Fatalf("expected 409, got %d", w.Code)
		}
	})
}

func TestVerification_Phone(t *testing.T) {
	db := testutil.SetupTestDB(t)
	router := testutil.NewRouter()
	t.Cleanup(func() { testutil.TruncateAll(db) })

	u := testutil.CreateUser(db, "phonecheck@example.com", "phoneuser")
	testutil.CreateProfile(db, u.UserID, "44999555555")
	token, _, _ := services_token(u.UserID, "client", "")

	t.Run("send phone verification", func(t *testing.T) {
		w := testutil.POST(router, "/api/v1/private/verify/phone/send").
			WithAuth(token).
			Do()

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d — body: %s", w.Code, w.Body.String())
		}

		data := testutil.DataField(w)
		msg, _ := data["message"].(string)
		if !containsStr(msg, "****") {
			t.Error("expected masked phone in message")
		}
	})

	t.Run("verify phone with correct code", func(t *testing.T) {
		var code string
		db.Raw("SELECT code FROM verifications WHERE user_id = ? AND type = 'phone' AND used_at IS NULL AND deleted_at IS NULL ORDER BY created_at DESC LIMIT 1", u.UserID).Scan(&code)

		w := testutil.POST(router, "/api/v1/private/verify/phone").
			WithAuth(token).
			WithBody(testutil.J("code", code)).
			Do()

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d — body: %s", w.Code, w.Body.String())
		}

		pw := testutil.GET(router, "/api/v1/private/profile/").WithAuth(token).Do()
		pdata := testutil.DataField(pw)
		if pdata["verified_phone"] != true {
			t.Error("profile.verified_phone should be true")
		}
	})
}

func containsStr(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		func() bool {
			for i := 0; i <= len(s)-len(substr); i++ {
				if s[i:i+len(substr)] == substr {
					return true
				}
			}
			return false
		}())
}
