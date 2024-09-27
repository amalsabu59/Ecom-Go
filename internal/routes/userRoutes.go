package routes

import (
	"gologin/internal/handlers"
	"gologin/internal/middleware"
	"net/http"
)

func UserRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/e", handlers.GetUserWithAddresses)
	mux.HandleFunc("/users/signup", handlers.SignUp)
	mux.HandleFunc("/users/login", handlers.Login)
	mux.Handle("/users/profile", middleware.Authenticate(http.HandlerFunc(handlers.GetUserWithAddresses)))
}
