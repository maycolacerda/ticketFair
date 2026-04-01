// tests/integration/purchase_test.go
package integration

import (
	"net/http"
	"testing"

	"github.com/maycolacerda/ticketfair/tests/testutil"
)

func TestPurchase_FullFlow(t *testing.T) {
	db := testutil.SetupTestDB(t)
	router := testutil.NewRouter()
	t.Cleanup(func() { testutil.TruncateAll(db) })

	// Seed
	m := testutil.CreateMerchant(db, "purchasemerchant@example.com")
	e := testutil.CreateEvent(db, m.MerchantID)
	tt := testutil.CreateTicketType(db, e.EventID, "general", 5000, 100)
	u := testutil.CreateUser(db, "buyer@example.com", "buyer")
	userToken, _, _ := services_token(u.UserID, "client", "")
	merchantToken, _, _ := services_token(m.MerchantID, "merchant", m.MerchantID)

	var transactionID string
	var ticketID string

	t.Run("purchase ticket", func(t *testing.T) {
		w := testutil.POST(router, "/api/v1/private/tickets/purchase").
			WithAuth(userToken).
			WithBody(testutil.J(
				"event_id", e.EventID,
				"ticket_type_id", tt.TicketTypeID,
				"quantity", 1,
			)).Do()

		if w.Code != http.StatusCreated {
			t.Fatalf("expected 201, got %d — body: %s", w.Code, w.Body.String())
		}

		data := testutil.DataField(w)
		transactionID, _ = data["transaction_id"].(string)
		if transactionID == "" {
			t.Fatal("expected transaction_id in response")
		}
		if data["status"] != "completed" {
			t.Errorf("expected status completed, got %v", data["status"])
		}
	})

	t.Run("ticket type availability decremented", func(t *testing.T) {
		var available int
		db.Raw("SELECT available FROM ticket_types WHERE ticket_type_id = ?", tt.TicketTypeID).Scan(&available)

		if available != 99 {
			t.Errorf("expected available 99, got %d", available)
		}
	})

	t.Run("ticket created for user", func(t *testing.T) {
		w := testutil.GET(router, "/api/v1/private/tickets/").
			WithAuth(userToken).
			Do()

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d — body: %s", w.Code, w.Body.String())
		}

		var resp struct {
			Data  []map[string]interface{} `json:"data"`
			Total float64                  `json:"total"`
		}
		testutil.ParseBody(w, &resp)

		if resp.Total < 1 {
			t.Fatal("expected at least 1 ticket")
		}

		ticketID, _ = resp.Data[0]["ticket_id"].(string)
		if resp.Data[0]["status"] != "active" {
			t.Errorf("expected status active, got %v", resp.Data[0]["status"])
		}
		if resp.Data[0]["ticket_type_name"] != tt.Name {
			t.Errorf("expected ticket_type_name %s, got %v", tt.Name, resp.Data[0]["ticket_type_name"])
		}
	})

	t.Run("get ticket by id", func(t *testing.T) {
		w := testutil.GET(router, "/api/v1/private/tickets/"+ticketID).
			WithAuth(userToken).
			Do()

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}

		data := testutil.DataField(w)
		if data["ticket_id"] != ticketID {
			t.Errorf("ticket_id mismatch")
		}
	})

	t.Run("validate ticket at door", func(t *testing.T) {
		w := testutil.POST(router, "/api/v1/merchant/tickets/"+ticketID+"/validate").
			WithAuth(merchantToken).
			Do()

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d — body: %s", w.Code, w.Body.String())
		}

		data := testutil.DataField(w)
		if data["status"] != "used" {
			t.Errorf("expected status used, got %v", data["status"])
		}
	})

	t.Run("cannot validate used ticket again", func(t *testing.T) {
		w := testutil.POST(router, "/api/v1/merchant/tickets/"+ticketID+"/validate").
			WithAuth(merchantToken).
			Do()

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected 400 for already-used ticket, got %d", w.Code)
		}
	})

	t.Run("cannot refund used ticket", func(t *testing.T) {
		w := testutil.POST(router, "/api/v1/private/tickets/refund").
			WithAuth(userToken).
			WithBody(testutil.J("transaction_id", transactionID)).
			Do()

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected 400 for refunding used ticket, got %d", w.Code)
		}
	})
}

func TestPurchase_SoldOut(t *testing.T) {
	db := testutil.SetupTestDB(t)
	router := testutil.NewRouter()
	t.Cleanup(func() { testutil.TruncateAll(db) })

	m := testutil.CreateMerchant(db, "soldout@example.com")
	e := testutil.CreateEvent(db, m.MerchantID)
	tt := testutil.CreateSoldOutTicketType(db, e.EventID)
	u := testutil.CreateUser(db, "buyer2@example.com", "buyer2")
	userToken, _, _ := services_token(u.UserID, "client", "")

	t.Run("purchase returns 409 when sold out", func(t *testing.T) {
		w := testutil.POST(router, "/api/v1/private/tickets/purchase").
			WithAuth(userToken).
			WithBody(testutil.J(
				"event_id", e.EventID,
				"ticket_type_id", tt.TicketTypeID,
				"quantity", 1,
			)).Do()

		if w.Code != http.StatusConflict {
			t.Fatalf("expected 409, got %d — body: %s", w.Code, w.Body.String())
		}
	})
}

func TestPurchase_SaleWindow(t *testing.T) {
	db := testutil.SetupTestDB(t)
	router := testutil.NewRouter()
	t.Cleanup(func() { testutil.TruncateAll(db) })

	m := testutil.CreateMerchant(db, "window@example.com")
	e := testutil.CreateEvent(db, m.MerchantID)
	u := testutil.CreateUser(db, "windowbuyer@example.com", "windowbuyer")
	userToken, _, _ := services_token(u.UserID, "client", "")

	t.Run("purchase within sale window succeeds", func(t *testing.T) {
		tt := testutil.CreateEarlyBirdTicketType(db, e.EventID)

		w := testutil.POST(router, "/api/v1/private/tickets/purchase").
			WithAuth(userToken).
			WithBody(testutil.J(
				"event_id", e.EventID,
				"ticket_type_id", tt.TicketTypeID,
				"quantity", 1,
			)).Do()

		if w.Code != http.StatusCreated {
			t.Fatalf("expected 201, got %d — body: %s", w.Code, w.Body.String())
		}
	})

	t.Run("purchase with quantity exceeding max returns 400", func(t *testing.T) {
		tt := testutil.CreateTicketType(db, e.EventID, "general", 5000, 100)
		// max_per_order defaults to 10, so 11 should fail
		db.Model(tt).Update("max_per_order", 5)

		w := testutil.POST(router, "/api/v1/private/tickets/purchase").
			WithAuth(userToken).
			WithBody(testutil.J(
				"event_id", e.EventID,
				"ticket_type_id", tt.TicketTypeID,
				"quantity", 6,
			)).Do()

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected 400 for exceeding max, got %d — body: %s", w.Code, w.Body.String())
		}
	})
}

func TestRefund_Flow(t *testing.T) {
	db := testutil.SetupTestDB(t)
	router := testutil.NewRouter()
	t.Cleanup(func() { testutil.TruncateAll(db) })

	m := testutil.CreateMerchant(db, "refundmerchant@example.com")
	e := testutil.CreateEventWithCapacity(db, m.MerchantID, 50)
	tt := testutil.CreateTicketType(db, e.EventID, "general", 5000, 50)
	u := testutil.CreateUser(db, "refundbuyer@example.com", "refundbuyer")
	userToken, _, _ := services_token(u.UserID, "client", "")

	t.Run("refund restores availability", func(t *testing.T) {
		// Purchase
		pw := testutil.POST(router, "/api/v1/private/tickets/purchase").
			WithAuth(userToken).
			WithBody(testutil.J(
				"event_id", e.EventID,
				"ticket_type_id", tt.TicketTypeID,
				"quantity", 1,
			)).Do()

		if pw.Code != http.StatusCreated {
			t.Fatalf("purchase failed: %d — %s", pw.Code, pw.Body.String())
		}

		txID := testutil.DataString(pw, "transaction_id")

		// Check availability before refund
		var availBefore int
		db.Raw("SELECT available FROM ticket_types WHERE ticket_type_id = ?", tt.TicketTypeID).Scan(&availBefore)

		// Refund
		rw := testutil.POST(router, "/api/v1/private/tickets/refund").
			WithAuth(userToken).
			WithBody(testutil.J("transaction_id", txID)).
			Do()

		if rw.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d — body: %s", rw.Code, rw.Body.String())
		}

		// Check availability restored
		var availAfter int
		db.Raw("SELECT available FROM ticket_types WHERE ticket_type_id = ?", tt.TicketTypeID).Scan(&availAfter)

		if availAfter != availBefore+1 {
			t.Errorf("expected available to be restored: before=%d after=%d", availBefore, availAfter)
		}
	})

	t.Run("cannot refund someone else's ticket", func(t *testing.T) {
		u2 := testutil.CreateUser(db, "intruder@example.com", "intruder")
		intruderToken, _, _ := services_token(u2.UserID, "client", "")

		// buyer purchases
		pw := testutil.POST(router, "/api/v1/private/tickets/purchase").
			WithAuth(userToken).
			WithBody(testutil.J(
				"event_id", e.EventID,
				"ticket_type_id", tt.TicketTypeID,
				"quantity", 1,
			)).Do()
		txID := testutil.DataString(pw, "transaction_id")

		// intruder tries to refund
		rw := testutil.POST(router, "/api/v1/private/tickets/refund").
			WithAuth(intruderToken).
			WithBody(testutil.J("transaction_id", txID)).
			Do()

		if rw.Code != http.StatusNotFound {
			t.Fatalf("expected 404 for unauthorized refund, got %d", rw.Code)
		}
	})
}

func TestTransactions_List(t *testing.T) {
	db := testutil.SetupTestDB(t)
	router := testutil.NewRouter()
	t.Cleanup(func() { testutil.TruncateAll(db) })

	m := testutil.CreateMerchant(db, "txmerchant@example.com")
	e := testutil.CreateEvent(db, m.MerchantID)
	tt := testutil.CreateTicketType(db, e.EventID, "general", 5000, 100)
	u := testutil.CreateUser(db, "txuser@example.com", "txuser")
	userToken, _, _ := services_token(u.UserID, "client", "")

	// Make 3 purchases
	for i := 0; i < 3; i++ {
		testutil.POST(router, "/api/v1/private/tickets/purchase").
			WithAuth(userToken).
			WithBody(testutil.J(
				"event_id", e.EventID,
				"ticket_type_id", tt.TicketTypeID,
				"quantity", 1,
			)).Do()
	}

	t.Run("returns paginated transactions for user", func(t *testing.T) {
		w := testutil.GET(router, "/api/v1/private/transactions/?page=1&limit=10").
			WithAuth(userToken).
			Do()

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}

		var resp struct {
			Data  []map[string]interface{} `json:"data"`
			Total float64                  `json:"total"`
		}
		testutil.ParseBody(w, &resp)

		if resp.Total < 3 {
			t.Errorf("expected at least 3 transactions, got %.0f", resp.Total)
		}
	})

	t.Run("no token returns 401", func(t *testing.T) {
		w := testutil.GET(router, "/api/v1/private/transactions/").Do()

		if w.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401, got %d", w.Code)
		}
	})
}
