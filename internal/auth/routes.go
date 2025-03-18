package auth

import (
	"chatrblox/internal/models"
	"chatrblox/pkg/utils"
	"encoding/json"
	"net/http"

	"gorm.io/gorm"
)

type AuthHandler struct {
	DB *gorm.DB
}

// User registration.
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	json.NewDecoder(r.Body).Decode(&input)

	hash, _ := utils.HashPassword(input.Password)

	user := models.User{
		Email:        input.Email,
		PasswordHash: hash,
	}

	if err := h.DB.Create(&user).Error; err != nil {
		http.Error(w, "User creation failed", http.StatusInternalServerError)
		return
	}

	token, _ := utils.GenerateJWT(user.ID)
	w.Write([]byte(token))
}

// Handle user login and authentication.
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	json.NewDecoder(r.Body).Decode(&input)

	var user models.User
	if err := h.DB.First(&user, "email = ?", input.Email).Error; err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	if !utils.CheckPassword(input.Password, user.PasswordHash) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, _ := utils.GenerateJWT(user.ID)
	w.Write([]byte(token))
}
