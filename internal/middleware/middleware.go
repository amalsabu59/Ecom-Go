package middleware

import (
	"context"
	"fmt"
	"gologin/internal/utils"
	"net/http"
)

type key string

const UserIDKey key = "userID"

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Retrieve the JWT token from cookies
		cookie, err := r.Cookie("token")
		if err != nil {
			utils.WriteResponse(w, http.StatusUnauthorized, "Unauthorized", nil, nil)
			return
		}

		tokenStr := cookie.Value

		// Validate the token and extract the claims
		claims, err := utils.ValidateJWT(tokenStr)
		if err != nil {
			utils.WriteResponse(w, http.StatusUnauthorized, "Token validation failed", nil, err)
			return
		}
		// Store the userID (or any other claims) in the context
		fmt.Println(claims.UserID)
		// ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		// Pass the updated request with the new context to the next handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RequireMethod(method string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			utils.WriteResponse(w, http.StatusMethodNotAllowed, "Invalid request method", nil, nil)
			return
		}
		next.ServeHTTP(w, r)
	})
}
