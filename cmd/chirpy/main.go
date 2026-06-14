package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/JorgeLR0610/Chirpy/internal/api"
	"github.com/JorgeLR0610/Chirpy/internal/repository"
	"github.com/JorgeLR0610/Chirpy/internal/service"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()
	platform := os.Getenv("PLATFORM")
	dbURL := os.Getenv("DB_URL")
	secret := os.Getenv("JWT_SECRET")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	const port = "8080"
	const filepathRoot = "web"

	// Repository
	dbQueries := repository.New(db)

	// Services
	adminSvc := service.NewAdminService(dbQueries)
	userSvc := service.NewUserService(dbQueries)
	chirpSvc := service.NewChirpService(dbQueries)
	tokenSvc := service.NewTokenService(dbQueries)

	// Handlers
	hitsCounter := &atomic.Int32{}

	adminHandler := api.NewAdminHandler(adminSvc, platform, hitsCounter)
	userHandler := api.NewUserHandler(userSvc, secret)
	chirpHandler := api.NewChirpHandler(chirpSvc, secret)
	tokensHandler := api.NewRefreshAccessTokenHandler(tokenSvc, secret)	

	mux := http.NewServeMux()
	
	// Endpoints

	// Static files server
	mux.Handle("/app/", adminHandler.MiddlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(filepathRoot)))))

	// Health check
	mux.HandleFunc("GET /api/healthz", api.HandlerReadiness)

	// Administration
	mux.HandleFunc("GET /admin/metrics", adminHandler.HandlerMetrics)
	mux.HandleFunc("POST /admin/reset", adminHandler.HandlerReset)

	// Users
	mux.HandleFunc("POST /api/users", userHandler.HandlerCreateUser)
	mux.HandleFunc("POST /api/login", userHandler.HandlerLoginUser)

	// Tokens
	mux.HandleFunc("POST /api/refresh", tokensHandler.HandlerRefreshAccessToken)
	mux.HandleFunc("POST /api/revoke", tokensHandler.HandlerRevokeRefreshToken)

	// Chirps
	mux.HandleFunc("POST /api/chirps", chirpHandler.HandlerCreateChirp)
	mux.HandleFunc("GET /api/chirps", chirpHandler.HandlerGetChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", chirpHandler.HandlerGetChirp)


	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
