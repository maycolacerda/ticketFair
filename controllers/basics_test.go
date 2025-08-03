package controllers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func setupTestRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/", GetHome)
	r.GET("/health", HealthCheck)
	r.POST("/users/new", NewUser)
	r.GET("/users", GetUsers)
	r.GET("/users/:id", GetUserByID)
	r.NoRoute(NotFound)
	return r
}

func TestGetHome(t *testing.T) {
	mockResponse := `{"message":"Welcome to the Ticket Fair API!"}`
	r := setupTestRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", w.Code)
	}
	if w.Body.String() != mockResponse {
		t.Errorf("Expected response body %q, got %q", mockResponse, w.Body.String())
	}
}

func TestHealthCheck(t *testing.T) {

	r := setupTestRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", w.Code)
	}
}
