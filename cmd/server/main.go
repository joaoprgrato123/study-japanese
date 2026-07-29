package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joaoprgrato123/study-japanese/internal/db"
	"github.com/joaoprgrato123/study-japanese/internal/middleware"
	"github.com/joaoprgrato123/study-japanese/internal/session"
	"github.com/joaoprgrato123/study-japanese/internal/user"
	"github.com/joho/godotenv"
)

func main() {

	//load .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal(" error loading .env file")

	}

	// connect to database
	ctx := context.Background()
	pool, err := db.Connect(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to database %v", err)
	}
	defer pool.Close()
	log.Println("connected to database")

	// repositories
	userRepo := user.NewUserRepository(pool)
	sessionRepo := session.NewSessionRepository(pool)

	// services
	userService := user.NewUserService(userRepo)
	sessionService := session.NewSessionService(sessionRepo)

	// handlers
	userHandler := user.NewUserHandler(userService, sessionService)

	// router

	r := mux.NewRouter()

	// public routes — no auth needed
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/register", userHandler.Register).Methods("POST")
	api.HandleFunc("/login", userHandler.Login).Methods("POST")
	api.HandleFunc("/logout", userHandler.Logout).Methods("POST")

	// protected routes — RequireAuth runs before every handler here
	protected := r.PathPrefix("/api").Subrouter()
	protected.Use(middleware.RequireAuth(sessionService))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("server running on http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("server error: %v", err)
	}

}
