package main

import (
	"gologin/internal/db"
	"gologin/internal/logger"
	"gologin/internal/routes"
	"net/http"
)

func main() {
	logger.SetupLogger()
	db.SetupDB()
	mux := http.NewServeMux()
	routes.UserRoutes(mux)
	logger.Log.Info().Msg("Server starting at :8080")
	http.ListenAndServe(":8080", mux)
}
