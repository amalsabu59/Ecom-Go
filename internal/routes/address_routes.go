package routes

import (
	"gologin/internal/handlers"
	"gologin/internal/middleware"
	"net/http"
)

func AddressRoutes(mux *http.ServeMux) {
	mux.Handle("/users/address", middleware.Authenticate(http.HandlerFunc(handlers.AddAddress)))
	mux.Handle("/users/addres/", middleware.Authenticate(middleware.RequireMethod("PUT", http.HandlerFunc(handlers.UpdateAddress))))
}
