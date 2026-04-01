// tests/integration/admin_test.go
package integration

import (
	"net/http"
	"testing"

	"github.com/maycolacerda/ticketfair/tests/testutil"
)

func TestAdmin_Auth(t *testing.T) {
	db := testutil.SetupTestDB(t)
	router := testutil.NewRouter()
	t.Cleanup(func() { testutil.TruncateAll(db) })

	testutil.CreateAdmin(db)

	t.Run("admin login success", func(t *testing.T) {
		w := testutil.POST(router, "/api/v1/admin/auth/login").
			WithBody(testutil.J(
				"email", "admin@test.com",
				"password", testutil.TestPassword,
			)).Do()

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d — body: %s", w.Code, w.Body.String())
		}

		data := testutil.DataField(w)
		if data["token"] == "" {
			t.Fatal("expected token")
		}
		if data["admin"] == nil {
			t.Fatal("expected admin object")
		}
	})

	t.Run("wrong password returns 401", func(t *testing.T) {
		w := testutil.POST(router, "/api/v1/admin/auth/login").
			WithBody(testutil.J(
				"email", "admin@test.com",
				"password", "Wrong!",
			)).Do()

		if w.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401, got %d", w.Code)
		}
	})
}

func TestAdmin_UserManagement(t *testing.T) {
	db := testutil.SetupTestDB(t)
	router := testutil.NewRouter()
	t.Cleanup(func() { testutil.TruncateAll(db) })

	admin := testutil.CreateAdmin(db)
	adminToken, _, _ := services_token(admin.AdminID, "superadmin", "")

	// Seed some users
	testutil.CreateUser(db, "user1@example.com", "user1")
	testutil.CreateUser(db, "user2@example.com", "user2")

	t.Run("list all users including inactive", func(t *testing.T) {
		u := testutil.CreateUser(db, "inactive@example.com", "inactiveuser")
		db.Model(u).Update("active", false)

		w := testutil.GET(router, "/api/v1/admin/users/").
			WithAuth(adminToken).
			Do()

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}

		var resp struct {
			Total float64 `json:"total"`
		}
		testutil.ParseBody(w, &resp)

		// Should include inactive user
		if resp.Total < 3 {
			t.Errorf("expected at least 3 users (incl. inactive), got %.0f", resp.Total)
		}
	})

	t.Run("admin creates user", func(t *testing.T) {
		w := testutil.POST(router, "/api/v1/admin/users/").
			WithAuth(adminToken).
			WithBody(testutil.J(
				"email", "adminmade@example.com",
				"username", "adminmade",
				"password", "PassW0rd!",
			)).Do()

		if w.Code != http.StatusCreated {
			t.Fatalf("expected 201, got %d — body: %s", w.Code, w.Body.String())
		}
	})

	t.Run("admin deactivates user", func(t *testing.T) {
		u := testutil.CreateUser(db, "todeactivate@example.com", "todeactivate")

		w := testutil.POST(router, "/api/v1/admin/users/"+u.UserID+"/deactivate").
			WithAuth(adminToken).
			Do()

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d — body: %s", w.Code, w.Body.String())
		}

		// User should not be able to login
		lw := testutil.POST(router, "/api/v1/public/auth/client/login").
			WithBody(testutil.J(
				"email", "todeactivate@example.com",
				"password", testutil.TestPassword,
			)).Do()

		if lw.Code != http.StatusForbidden {
			t.Errorf("expected 403 after deactivation, got %d", lw.Code)
		}
	})

	t.Run("admin activates user", func(t *testing.T) {
		u := testutil.CreateUser(db, "toreactivate@example.com", "toreactivate")
		db.Model(u).Update("active", false)

		w := testutil.POST(router, "/api/v1/admin/users/"+u.UserID+"/activate").
			WithAuth(adminToken).
			Do()

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}

		// Should be able to login now
		lw := testutil.POST(router, "/api/v1/public/auth/client/login").
			WithBody(testutil.J(
				"email", "toreactivate@example.com",
				"password", testutil.TestPassword,
			)).Do()

		if lw.Code != http.StatusOK {
			t.Errorf("expected 200 after activation, got %d", lw.Code)
		}
	})

	t.Run("client token cannot access admin routes", func(t *testing.T) {
		u := testutil.CreateUser(db, "client@example.com", "clientonly")
		clientToken, _, _ := services_token(u.UserID, "client", "")

		w := testutil.GET(router, "/api/v1/admin/users/").
			WithAuth(clientToken).
			Do()

		if w.Code != http.StatusForbidden {
			t.Fatalf("expected 403, got %d", w.Code)
		}
	})

	t.Run("no token returns 401", func(t *testing.T) {
		w := testutil.GET(router, "/api/v1/admin/users/").Do()

		if w.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401, got %d", w.Code)
		}
	})
}

func TestAdmin_MerchantManagement(t *testing.T) {
	db := testutil.SetupTestDB(t)
	router := testutil.NewRouter()
	t.Cleanup(func() { testutil.TruncateAll(db) })

	admin := testutil.CreateAdmin(db)
	adminToken, _, _ := services_token(admin.AdminID, "superadmin", "")
	m := testutil.CreateMerchant(db, "managed@example.com")
	r := testutil.CreateMerchantRep(db, m.MerchantID, "managedRep@example.com", "admin")

	t.Run("list all merchants", func(t *testing.T) {
		w := testutil.GET(router, "/api/v1/admin/merchants/").
			WithAuth(adminToken).
			Do()

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}
	})

	t.Run("admin creates merchant", func(t *testing.T) {
		w := testutil.POST(router, "/api/v1/admin/merchants/").
			WithAuth(adminToken).
			WithBody(testutil.J(
				"name", "New Merchant",
				"email", "newmerchant@example.com",
				"password", "PassW0rd!",
				"phone", "44988000000",
				"description", "Created by admin",
			)).Do()

		if w.Code != http.StatusCreated {
			t.Fatalf("expected 201, got %d — body: %s", w.Code, w.Body.String())
		}
	})

	t.Run("deactivate merchant blocks reps", func(t *testing.T) {
		// Rep should be able to login before
		before := testutil.POST(router, "/api/v1/public/auth/rep/login").
			WithBody(testutil.J(
				"email", r.Email,
				"password", testutil.TestPassword,
			)).Do()

		if before.Code != http.StatusOK {
			t.Fatalf("rep login should succeed before deactivation, got %d", before.Code)
		}

		// Deactivate merchant
		dw := testutil.POST(router, "/api/v1/admin/merchants/"+m.MerchantID+"/deactivate").
			WithAuth(adminToken).
			Do()

		if dw.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d — body: %s", dw.Code, dw.Body.String())
		}

		// Rep login should now fail
		after := testutil.POST(router, "/api/v1/public/auth/rep/login").
			WithBody(testutil.J(
				"email", r.Email,
				"password", testutil.TestPassword,
			)).Do()

		if after.Code != http.StatusForbidden {
			t.Errorf("rep login should fail after merchant deactivation, got %d", after.Code)
		}
	})

	t.Run("activate merchant restores rep access", func(t *testing.T) {
		testutil.POST(router, "/api/v1/admin/merchants/"+m.MerchantID+"/activate").
			WithAuth(adminToken).
			Do()

		w := testutil.POST(router, "/api/v1/public/auth/rep/login").
			WithBody(testutil.J(
				"email", r.Email,
				"password", testutil.TestPassword,
			)).Do()

		if w.Code != http.StatusOK {
			t.Errorf("rep login should succeed after merchant activation, got %d", w.Code)
		}
	})

	t.Run("duplicate merchant email returns 409", func(t *testing.T) {
		w := testutil.POST(router, "/api/v1/admin/merchants/").
			WithAuth(adminToken).
			WithBody(testutil.J(
				"name", "Dup Merchant",
				"email", "managed@example.com", // already exists
				"password", "PassW0rd!",
				"phone", "44988000001",
			)).Do()

		if w.Code != http.StatusConflict {
			t.Fatalf("expected 409, got %d", w.Code)
		}
	})

	// Silence unused variable warning
	_ = r
}
