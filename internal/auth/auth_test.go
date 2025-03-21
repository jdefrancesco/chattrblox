package auth

import (
	"bytes"
	"chatrblox/internal/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to test DB: %v", err)
	}
	db.AutoMigrate(&models.User{})
	return db
}

func TestRegister(t *testing.T) {
	db := setupTestDB(t)
	handler := &AuthHandler{DB: db}

	body := map[string]string{
		"email":    "test@example.com",
		"password": "testpass",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.Register(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}

	var user models.User
	err := db.First(&user, "email = ?", body["email"]).Error
	if err != nil {
		t.Errorf("user not created in DB")
	}
}

func TestLogin(t *testing.T) {
	db := setupTestDB(t)
	handler := &AuthHandler{DB: db}

	// Manually create user
	user := models.User{
		Email:        "login@example.com",
		PasswordHash: "$2a$14$RpS8Zta3tZptQqpFtGeIvuJh0i2d6gLQF3rptCrRfhYydgic4N1rS", // "secret"
	}
	db.Create(&user)

	body := map[string]string{
		"email":    "login@example.com",
		"password": "secret",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.Login(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}

	if rr.Body.Len() == 0 {
		t.Error("expected token in response body")
	}
}
