// tests/integration/event_test.go
package integration

import (
	"net/http"
	"testing"
	"time"

	"github.com/maycolacerda/ticketfair/tests/testutil"
)

func TestEvents_Public(t *testing.T) {
	db := testutil.SetupTestDB(t)
	router := testutil.NewRouter()
	t.Cleanup(func() { testutil.TruncateAll(db) })

	m := testutil.CreateMerchant(db, "eventmerchant@example.com")
	testutil.CreateEvent(db, m.MerchantID)
	testutil.CreateEvent(db, m.MerchantID)

	t.Run("list events returns active upcoming events", func(t *testing.T) {
		w := testutil.GET(router, "/api/v1/public/events/").Do()

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d — body: %s", w.Code, w.Body.String())
		}

		var resp struct {
			Data  []map[string]interface{} `json:"data"`
			Total float64                  `json:"total"`
		}
		testutil.ParseBody(w, &resp)

		if resp.Total < 2 {
			t.Errorf("expected at least 2 events, got %.0f", resp.Total)
		}
	})

	t.Run("get event by id", func(t *testing.T) {
		e := testutil.CreateEvent(db, m.MerchantID)

		w := testutil.GET(router, "/api/v1/public/events/"+e.EventID).Do()

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}

		data := testutil.DataField(w)
		if data["event_id"] != e.EventID {
			t.Errorf("expected event_id %s, got %v", e.EventID, data["event_id"])
		}
	})

	t.Run("get event by invalid id returns 404", func(t *testing.T) {
		w := testutil.GET(router, "/api/v1/public/events/00000000-0000-0000-0000-000000000000").Do()

		if w.Code != http.StatusNotFound {
			t.Fatalf("expected 404, got %d", w.Code)
		}
	})
}

func TestEvents_Merchant(t *testing.T) {
	db := testutil.SetupTestDB(t)
	router := testutil.NewRouter()
	t.Cleanup(func() { testutil.TruncateAll(db) })

	m := testutil.CreateMerchant(db, "evmerchant@example.com")
	token, _, _ := services_token(m.MerchantID, "merchant", m.MerchantID)

	t.Run("create event", func(t *testing.T) {
		w := testutil.POST(router, "/api/v1/merchant/events/new").
			WithAuth(token).
			WithBody(testutil.J(
				"name", "Integration Test Event",
				"description", "A test event",
				"location", "Test Venue",
				"start_time", time.Now().Add(30*24*time.Hour).Format(time.RFC3339),
				"end_time", time.Now().Add(30*24*time.Hour+4*time.Hour).Format(time.RFC3339),
				"capacity", 200,
			)).Do()

		if w.Code != http.StatusCreated {
			t.Fatalf("expected 201, got %d — body: %s", w.Code, w.Body.String())
		}

		data := testutil.DataField(w)
		if data["event_id"] == "" {
			t.Error("expected event_id in response")
		}
		if data["merchant_id"] != m.MerchantID {
			t.Errorf("expected merchant_id %s, got %v", m.MerchantID, data["merchant_id"])
		}
	})

	t.Run("create event with past date returns 400", func(t *testing.T) {
		w := testutil.POST(router, "/api/v1/merchant/events/new").
			WithAuth(token).
			WithBody(testutil.J(
				"name", "Past Event",
				"description", "Already happened",
				"location", "Somewhere",
				"start_time", time.Now().Add(-24*time.Hour).Format(time.RFC3339),
				"end_time", time.Now().Add(-20*time.Hour).Format(time.RFC3339),
				"capacity", 100,
			)).Do()

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d — body: %s", w.Code, w.Body.String())
		}
	})

	t.Run("update event", func(t *testing.T) {
		e := testutil.CreateEvent(db, m.MerchantID)

		w := testutil.PUT(router, "/api/v1/merchant/events/"+e.EventID).
			WithAuth(token).
			WithBody(testutil.J(
				"name", "Updated Event Name",
				"capacity", 300,
			)).Do()

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d — body: %s", w.Code, w.Body.String())
		}

		data := testutil.DataField(w)
		if data["name"] != "Updated Event Name" {
			t.Errorf("expected updated name, got %v", data["name"])
		}
		if data["capacity"] != float64(300) {
			t.Errorf("expected capacity 300, got %v", data["capacity"])
		}
	})

	t.Run("cannot update another merchant's event", func(t *testing.T) {
		m2 := testutil.CreateMerchant(db, "other@example.com")
		e2 := testutil.CreateEvent(db, m2.MerchantID)

		w := testutil.PUT(router, "/api/v1/merchant/events/"+e2.EventID).
			WithAuth(token).
			WithBody(testutil.J("name", "Hijacked")).
			Do()

		if w.Code != http.StatusNotFound {
			t.Fatalf("expected 404 when updating another merchant's event, got %d", w.Code)
		}
	})

	t.Run("no token returns 401", func(t *testing.T) {
		w := testutil.POST(router, "/api/v1/merchant/events/new").
			WithBody(testutil.J("name", "No Auth")).
			Do()

		if w.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401, got %d", w.Code)
		}
	})
}

func TestTicketTypes(t *testing.T) {
	db := testutil.SetupTestDB(t)
	router := testutil.NewRouter()
	t.Cleanup(func() { testutil.TruncateAll(db) })

	m := testutil.CreateMerchant(db, "ttmerchant@example.com")
	e := testutil.CreateEvent(db, m.MerchantID)
	token, _, _ := services_token(m.MerchantID, "merchant", m.MerchantID)

	t.Run("create general ticket type", func(t *testing.T) {
		w := testutil.POST(router, "/api/v1/merchant/events/"+e.EventID+"/ticket-types").
			WithAuth(token).
			WithBody(testutil.J(
				"name", "General Admission",
				"category", "general",
				"price_cents", 5000,
				"capacity", 200,
				"min_per_order", 1,
				"max_per_order", 10,
			)).Do()

		if w.Code != http.StatusCreated {
			t.Fatalf("expected 201, got %d — body: %s", w.Code, w.Body.String())
		}

		data := testutil.DataField(w)
		if data["ticket_type_id"] == "" {
			t.Error("expected ticket_type_id")
		}
		if data["price_formatted"] != "R$ 50,00" {
			t.Errorf("expected R$ 50,00 got %v", data["price_formatted"])
		}
		if data["available"] != float64(200) {
			t.Errorf("expected available 200, got %v", data["available"])
		}
		if data["on_sale"] != true {
			t.Error("expected on_sale = true")
		}
	})

	t.Run("create vip ticket type", func(t *testing.T) {
		w := testutil.POST(router, "/api/v1/merchant/events/"+e.EventID+"/ticket-types").
			WithAuth(token).
			WithBody(testutil.J(
				"name", "VIP",
				"category", "vip",
				"price_cents", 15000,
				"capacity", 50,
				"min_per_order", 1,
				"max_per_order", 2,
			)).Do()

		if w.Code != http.StatusCreated {
			t.Fatalf("expected 201, got %d — body: %s", w.Code, w.Body.String())
		}
	})

	t.Run("invalid category returns 422", func(t *testing.T) {
		w := testutil.POST(router, "/api/v1/merchant/events/"+e.EventID+"/ticket-types").
			WithAuth(token).
			WithBody(testutil.J(
				"name", "Invalid",
				"category", "platinum_ultra", // not in enum
				"price_cents", 9999,
				"capacity", 10,
			)).Do()

		if w.Code != http.StatusUnprocessableEntity {
			t.Fatalf("expected 422, got %d", w.Code)
		}
	})

	t.Run("public endpoint returns only active types", func(t *testing.T) {
		// Deactivate one type
		var ttID string
		db.Raw("SELECT ticket_type_id FROM ticket_types WHERE event_id = ? LIMIT 1", e.EventID).Scan(&ttID)

		testutil.PUT(router, "/api/v1/merchant/events/"+e.EventID+"/ticket-types/"+ttID).
			WithAuth(token).
			WithBody(testutil.J("active", false)).
			Do()

		w := testutil.GET(router, "/api/v1/public/events/"+e.EventID+"/ticket-types").Do()

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}

		var resp struct {
			Data []map[string]interface{} `json:"data"`
		}
		testutil.ParseBody(w, &resp)
		for _, tt := range resp.Data {
			if tt["active"] == false {
				t.Error("public endpoint should not return inactive ticket types")
			}
		}
	})

	t.Run("update ticket type adjusts available proportionally", func(t *testing.T) {
		tt := testutil.CreateTicketType(db, e.EventID, "general", 3000, 100)

		w := testutil.PUT(router, "/api/v1/merchant/events/"+e.EventID+"/ticket-types/"+tt.TicketTypeID).
			WithAuth(token).
			WithBody(testutil.J("capacity", 150)).
			Do()

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d — body: %s", w.Code, w.Body.String())
		}

		data := testutil.DataField(w)
		if data["capacity"] != float64(150) {
			t.Errorf("expected capacity 150, got %v", data["capacity"])
		}
		// available should increase by 50 (100 → 150, since none sold)
		if data["available"] != float64(150) {
			t.Errorf("expected available 150, got %v", data["available"])
		}
	})

	t.Run("delete ticket type with no sales", func(t *testing.T) {
		tt := testutil.CreateTicketType(db, e.EventID, "complimentary", 0, 20)

		w := testutil.DELETE(router, "/api/v1/merchant/events/"+e.EventID+"/ticket-types/"+tt.TicketTypeID).
			WithAuth(token).
			Do()

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d — body: %s", w.Code, w.Body.String())
		}
	})

	t.Run("delete ticket type with sales returns 400", func(t *testing.T) {
		tt := testutil.CreateTicketType(db, e.EventID, "general", 5000, 100)
		// Simulate a sale by decrementing available
		db.Model(tt).Update("available", 99)

		w := testutil.DELETE(router, "/api/v1/merchant/events/"+e.EventID+"/ticket-types/"+tt.TicketTypeID).
			WithAuth(token).
			Do()

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d — body: %s", w.Code, w.Body.String())
		}
	})
}

// thin wrapper — avoids circular import by using services via testutil.SetTokenGenerator
func services_token(userID, role, merchantID string) (string, int64, error) {
	return services_GenerateToken(userID, role, merchantID)
}
