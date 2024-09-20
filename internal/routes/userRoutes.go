package routes

import (
	"gologin/internal/handlers"
	"net/http"
)

func UserRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/users/sign_up", handlers.SignUp)
	mux.HandleFunc("/users/login", handlers.Login)
}
