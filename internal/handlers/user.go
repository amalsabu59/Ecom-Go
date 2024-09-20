package handlers

import (
	"encoding/json"
	"gologin/internal/db"
	"gologin/internal/logger"
	"gologin/internal/models"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func SignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if user.Name == "" || user.Email == "" || user.Password == "" {
		http.Error(w, "Name, Email, and Password are required", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to hash password")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)

	if err := ensureUsersTableExists(); err != nil {
		logger.Log.Error().Err(err).Msg("Failed to ensure users table exists")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	_, err = db.DB.NewInsert().Model(&user).Exec(r.Context())
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to insert user")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var loginDetails struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&loginDetails); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if loginDetails.Email == "" || loginDetails.Password == "" {
		http.Error(w, "Email and Password are required", http.StatusBadRequest)
		return
	}

	var user models.User
	err := db.DB.NewSelect().Model(&user).Where("email = ?", loginDetails.Email).Limit(1).Scan(r.Context())
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to find user")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginDetails.Password)); err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

		response := struct {
		Message string      `json:"message"`
		User    models.User `json:"user"`
	}{
		Message: "Logged in successfully",
		User:    user,
	}
	response.User.Password = "" 

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}


func ensureUsersTableExists() error {
	_, err := db.DB.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            id SERIAL PRIMARY KEY,
            name VARCHAR(100) NOT NULL,
            email VARCHAR(100) NOT NULL UNIQUE,
            password VARCHAR(255) NOT NULL
        );
    `)
	return err
}
