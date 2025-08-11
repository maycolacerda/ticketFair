package controllers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/maycolacerda/ticketfair/models"
	"gorm.io/gorm"
)

var UserTest = []struct {
	Name           string
	User           models.User
	WantErr        bool
	ExpectedStatus int
	ExpectedError  string
	MockCreate     func(value interface{}) *gorm.DB
}{
	{
		Name:           "Invalid Password",
		User:           models.User{Username: "TestUser", Email: "test@example.com", Password: "passW0rd@"},
		WantErr:        true,
		ExpectedStatus: http.StatusNotAcceptable,
		ExpectedError:  "Password must contain at least one special character",
	},
	{
		Name:           "Invalid Email",
		User:           models.User{Username: "TestUser", Email: "test", Password: "passW0rd@"},
		WantErr:        true,
		ExpectedStatus: http.StatusNotAcceptable,
		ExpectedError:  "Invalid email format",
	},
	{
		Name:           "Empty Username",
		User:           models.User{Username: "", Email: "test@example.com", Password: "passW0rd@"},
		WantErr:        true,
		ExpectedStatus: http.StatusNotAcceptable,
		ExpectedError:  "Username is required",
	},
	{
		Name:           "Empty Email",
		User:           models.User{Username: "TestUser", Email: "", Password: "passW0rd@"},
		WantErr:        true,
		ExpectedStatus: http.StatusNotAcceptable,
		ExpectedError:  "Email is required",
	},
	{
		Name:           "Empty Password",
		User:           models.User{Username: "TestUser", Email: "test@example.com", Password: ""},
		WantErr:        true,
		ExpectedStatus: http.StatusNotAcceptable,
		ExpectedError:  "Password is required",
	},
	{
		Name:           "Short Password",
		User:           models.User{Username: "TestUser", Email: "test@example.com", Password: "short"},
		WantErr:        true,
		ExpectedStatus: http.StatusNotAcceptable,
		ExpectedError:  "Password must be at least 8 characters long",
	},
	{
		Name:           "Missing Symbol",
		User:           models.User{Username: "TestUser", Email: "test@example.com", Password: "passW0rd"},
		WantErr:        true,
		ExpectedStatus: http.StatusNotAcceptable,
		ExpectedError:  "Password must contain at least one special character",
	},
	{
		Name:           "invalid chacter in username",
		User:           models.User{Username: "Test@User", Email: "test@example.com", Password: "passW0rd@"},
		WantErr:        true,
		ExpectedStatus: http.StatusNotAcceptable,
		ExpectedError:  "Username must not contain only letters and numbers",
	},
}

func TestNewUser(t *testing.T) {
	router := setupTestRouter()

	for _, tt := range UserTest {
		t.Run(tt.Name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/users/new", nil)
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)
			if (w.Code != http.StatusOK) != tt.WantErr {
				t.Errorf("Expected status code %d, got %d", tt.ExpectedStatus, w.Code)
			}
		})
	}
}
