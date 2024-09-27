package main

import (
	"gologin/internal/db"
	"gologin/internal/logger"
	"gologin/internal/routes"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		logger.Log.Error().Err(err).Msg("Error loading .env file")
	}
	logger.SetupLogger()
	db.SetupDB()
	// db.Migrate()
	mux := http.NewServeMux()
	routes.UserRoutes(mux)
	routes.AddressRoutes(mux)

	if len(os.Args) < 2 {
		logger.Log.Warn().Msg("Please provide a command: migrate or rollback")
	} else {
		switch os.Args[1] {
		case "migrate":
			if err := db.Migrate(); err != nil {
				logger.Log.Fatal().Err(err).Msg("Migration failed")
			}
			logger.Log.Info().Msg("Migration completed successfully")
			return // Exit after migration

		case "rollback":
			if err := db.Rollback(); err != nil {
				logger.Log.Fatal().Err(err).Msg("Rollback failed")
			}
			logger.Log.Info().Msg("Rollback completed successfully")
			return // Exit after rollback

		default:
			logger.Log.Fatal().Msg("Unknown command")
		}
	}

	cors := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),                                       // Allow all origins (you may want to restrict this in production)
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}), // Allow specific methods
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),           // Allow specific headers
	)

	server := &http.Server{
		Addr:         ":8081",
		Handler:      cors(mux),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	logger.Log.Info().Msg("Server starting at :8081")
	if err := server.ListenAndServe(); err != nil {
		logger.Log.Error().Err(err).Msg("Failed to start server")
	}

}
