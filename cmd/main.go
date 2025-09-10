package main

import (
	"chatrblox/internal/auth"
	"chatrblox/internal/middleware"
	"chatrblox/internal/models"
	"chatrblox/internal/ws"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	chiMid "github.com/go-chi/chi/v5/middleware"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Grab DB (Postgres) URL
	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		log.Fatalln("Set DATABASE_URL environment variable.")
	}
	// Get port from environment. Default to 8080 is missing.
	port := os.Getenv("PORT")
	if port == "" {
		log.Println("Setting to default port 8080")
		port = "8080"
	}

	db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}

	// Auto-migrate the user model
	if err := db.AutoMigrate(&models.User{}, &models.Report{}, &models.Session{}); err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}

	// Setup Redis. Matchmaking and other subsystems rely on this.
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// Match making subsystem init.
	matchmaker := ws.NewMatchmaker(redisClient)
	matchmaker.DB = db
	hub := &ws.Hub{Matchmaker: matchmaker}

	// Setup chi router.
	r := chi.NewRouter()
	r.Use(chiMid.Logger)

	authHandler := &auth.AuthHandler{DB: db}

	// Endpoints for registering and login.
	r.Post("/api/register", authHandler.Register)
	r.Post("/api/login", authHandler.Login)

	// Protected API routes. Using to test JWT are working.
	r.Group(func(r chi.Router) {
		r.Use(middleware.JWTAuth)
		r.Get("/api/protected", func(w http.ResponseWriter, r *http.Request) {
			userID := middleware.GetUserID(r)
			w.Write([]byte("Hello, user: " + userID.String()))
		})
	})

	// Middleware. Handles Authentication, etc..
	r.Group(func(r chi.Router) {
		r.Use(middleware.JWTAuth)
		r.Get("/ws", hub.HandleWS)
	})

	fmt.Printf("[*] Started server. Listening on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
