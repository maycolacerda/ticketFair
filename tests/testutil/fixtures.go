// tests/testutil/fixtures.go
package testutil

import (
	"fmt"
	"time"

	"github.com/maycolacerda/ticketfair/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const TestPassword = "PassW0rd!"

func hashPassword(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash)
}

// ── Users ─────────────────────────────────────────────

func CreateUser(db *gorm.DB, email, username string) *models.User {
	u := &models.User{
		Email:    email,
		Username: username,
		Password: hashPassword(TestPassword),
		Active:   true,
	}
	db.Create(u)
	return u
}

func CreateAdmin(db *gorm.DB) *models.Admin {
	a := &models.Admin{
		Name:     "Test Admin",
		Email:    "admin@test.com",
		Password: hashPassword(TestPassword),
		Active:   true,
	}
	db.Create(a)
	return a
}

// ── Merchants ─────────────────────────────────────────

func CreateMerchant(db *gorm.DB, email string) *models.Merchant {
	m := &models.Merchant{
		Name:        "Test Merchant",
		Email:       email,
		Password:    hashPassword(TestPassword),
		Phone:       "44999000000",
		Description: "Integration test merchant",
		Active:      true,
	}
	db.Create(m)
	return m
}

func CreateMerchantRep(db *gorm.DB, merchantID, email, role string) *models.MerchantRep {
	r := &models.MerchantRep{
		MerchantID: merchantID,
		Name:       "Test Rep",
		Email:      email,
		Password:   hashPassword(TestPassword),
		Phone:      "44999000001",
		Role:       role,
		Active:     true,
	}
	db.Create(r)
	return r
}

// ── Events ────────────────────────────────────────────

func CreateEvent(db *gorm.DB, merchantID string) *models.Event {
	e := &models.Event{
		MerchantID:  merchantID,
		Name:        "Test Event",
		Description: "Integration test event",
		Location:    "Test Venue",
		StartTime:   time.Now().Add(30 * 24 * time.Hour),
		EndTime:     time.Now().Add(30*24*time.Hour + 4*time.Hour),
		Capacity:    500,
		Active:      true,
	}
	db.Create(e)
	return e
}

func CreateEventWithCapacity(db *gorm.DB, merchantID string, capacity int) *models.Event {
	e := &models.Event{
		MerchantID: merchantID,
		Name:       fmt.Sprintf("Event (cap %d)", capacity),
		Location:   "Test Venue",
		StartTime:  time.Now().Add(30 * 24 * time.Hour),
		EndTime:    time.Now().Add(30*24*time.Hour + 4*time.Hour),
		Capacity:   capacity,
		Active:     true,
	}
	db.Create(e)
	return e
}

// ── Ticket Types ──────────────────────────────────────

func CreateTicketType(db *gorm.DB, eventID string, category string, priceCents int64, capacity int) *models.TicketType {
	tt := &models.TicketType{
		EventID:     eventID,
		Name:        fmt.Sprintf("%s Ticket", category),
		Category:    category,
		PriceCents:  priceCents,
		Capacity:    capacity,
		Available:   capacity,
		MinPerOrder: 1,
		MaxPerOrder: 10,
		Active:      true,
	}
	db.Create(tt)
	return tt
}

func CreateSoldOutTicketType(db *gorm.DB, eventID string) *models.TicketType {
	tt := &models.TicketType{
		EventID:    eventID,
		Name:       "Sold Out Ticket",
		Category:   "general",
		PriceCents: 5000,
		Capacity:   10,
		Available:  0, // fully sold
		Active:     true,
	}
	db.Create(tt)
	return tt
}

func CreateEarlyBirdTicketType(db *gorm.DB, eventID string) *models.TicketType {
	saleEnd := time.Now().Add(7 * 24 * time.Hour)
	tt := &models.TicketType{
		EventID:     eventID,
		Name:        "Early Bird",
		Category:    "early_bird",
		PriceCents:  2500,
		Capacity:    100,
		Available:   100,
		MinPerOrder: 1,
		MaxPerOrder: 4,
		SaleEndsAt:  &saleEnd,
		Active:      true,
	}
	db.Create(tt)
	return tt
}

// ── Profiles ──────────────────────────────────────────

func CreateProfile(db *gorm.DB, userID, phone string) *models.Profile {
	p := &models.Profile{
		UserID:      userID,
		FirstName:   "Test",
		LastName:    "User",
		PhoneNumber: phone,
	}
	db.Create(p)

	a := &models.Address{
		ProfileID: p.ProfileID,
		Street:    "Rua Test, 123",
		City:      "Cianorte",
		State:     "PR",
		Country:   "BR",
		ZipCode:   "87200000",
	}
	db.Create(a)
	p.Address = *a
	return p
}

// ── Tokens ────────────────────────────────────────────

func TokenForUser(userID string) string {
	token, _, _ := services_GenerateToken(userID, "client", "")
	return token
}

func TokenForMerchant(merchantID string) string {
	token, _, _ := services_GenerateToken(merchantID, "merchant", merchantID)
	return token
}

func TokenForAdmin(adminID string) string {
	token, _, _ := services_GenerateToken(adminID, "superadmin", "")
	return token
}

// services_GenerateToken is a thin wrapper to avoid importing the full
// services package and triggering DB init. Call services.GenerateToken directly
// in test files that already import services.
var services_GenerateToken func(userID, role, merchantID string) (string, int64, error)

func SetTokenGenerator(fn func(string, string, string) (string, int64, error)) {
	services_GenerateToken = fn
}
