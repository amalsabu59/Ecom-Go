package main

import (
	"gologin/internal/db"
	"gologin/internal/logger"
	"gologin/internal/routes"
	"net/http"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		logger.Log.Error().Err(err).Msg("Error loading .env file")
	}
	logger.SetupLogger()
	db.SetupDB()
	db.Migrate()
	mux := http.NewServeMux()
	routes.UserRoutes(mux)
	routes.AddressRoutes(mux)

	server := &http.Server{
		Addr:         ":8081",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	logger.Log.Info().Msg("Server starting at :8081")
	if err := server.ListenAndServe(); err != nil {
		logger.Log.Error().Err(err).Msg("Failed to start server")
	}
}
