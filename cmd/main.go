package main

import (
	"chatrblox/internal/auth"
	"chatrblox/internal/models"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dbUrl := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}

	// Auto-migrate the user model
	if err := db.AutoMigrate(&models.User{}); err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	authHandler := &auth.AuthHandler{DB: db}
	r.Post("/api/register", authHandler.Register)
	r.Post("/api/login", authHandler.Login)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("[*] Started server. Listening on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
