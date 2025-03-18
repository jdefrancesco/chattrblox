package main

import (
	"chatrblox/internal/auth"
	"chatrblox/internal/middleware"
	"chatrblox/internal/models"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	chiMid "github.com/go-chi/chi/v5/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Grab postgresql URL
	dbUrl := os.Getenv("DATABASE_URL")

	db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}

	// Auto-migrate the user model
	if err := db.AutoMigrate(&models.User{}, &models.Report{}, &models.Session{}); err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}

	// Setup chi router.
	r := chi.NewRouter()
	r.Use(chiMid.Logger)

	authHandler := &auth.AuthHandler{DB: db}
	r.Post("/api/register", authHandler.Register)
	r.Post("/api/login", authHandler.Login)

	// Protected API routes.
	r.Group(func(r chi.Router) {
		r.Use(middleware.JWTAuth)
		r.Get("/api/protected", func(w http.ResponseWriter, r *http.Request) {
			userID := middleware.GetUserID(r)
			w.Write([]byte("Hello, user: " + userID.String()))
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("[*] Started server. Listening on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
