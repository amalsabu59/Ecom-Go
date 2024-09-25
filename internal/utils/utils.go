// internal/utils/utils.go
package utils

import (
	"encoding/json"
	"gologin/internal/logger"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Secret key used to sign the JWT. Change it to something secure!
var jwtSecret = []byte(os.Getenv("JWT_SECRET_KEY")) // Change this to a secure secret

// Claims represents the JWT claims (payload)
type Claims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// ApiResponse represents a standard API response
type APIResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// GenerateJWT generates a new JWT token for a given username
func GenerateJWT(username string, userID int64) (string, error) {
	// Set expiration time for the token
	expirationTime := time.Now().Add(24 * time.Hour) // Token valid for 24 hours

	// Create the JWT claims, which includes the username and expiry time
	claims := &Claims{
		Username: username,
		UserID:   userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// Create the token with the HS256 signing algorithm
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	return token.SignedString(jwtSecret)
}

// ValidateJWT validates the provided JWT token string and returns the claims if valid
func ValidateJWT(tokenString string) (*Claims, error) {
	// Parse the token with the claims
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(_ *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	// Check if there is an error or if the token is invalid
	if err != nil || !token.Valid {
		logger.Log.Error().Err(err).Msg("Failed to Validate JWT Token")
		return nil, err
	}

	// Return the claims (username, etc.) if the token is valid
	return claims, nil
}

func WriteResponse(w http.ResponseWriter, statusCode int, message string, data interface{}, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	response := APIResponse{
		Message: message,
		Data:    data,
	}
	if err != nil {
		response.Error = err.Error()
	}

	if encodeErr := json.NewEncoder(w).Encode(response); encodeErr != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
