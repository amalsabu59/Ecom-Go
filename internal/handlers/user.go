package handlers

import (
	"encoding/json"
	"fmt"
	"gologin/internal/db"
	"gologin/internal/logger"
	"gologin/internal/models"
	"gologin/internal/utils"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func SignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteResponse(w, http.StatusMethodNotAllowed, "Invalid request method", nil, nil)
		return
	}

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		utils.WriteResponse(w, http.StatusBadRequest, "Invalid request payload", nil, err)
		return
	}

	// Check for required fields
	if user.Name == "" {
		utils.WriteResponse(w, http.StatusBadRequest, "Name is required", nil, nil)
		return
	}
	if user.Email == "" {
		utils.WriteResponse(w, http.StatusBadRequest, "Email is required", nil, nil)
		return
	}
	if user.Password == "" {
		utils.WriteResponse(w, http.StatusBadRequest, "Password is required", nil, nil)
		return
	}

	// Check if the email already exists
	var existingUser models.User
	err := db.DB.NewSelect().Model(&existingUser).Where("email = ?", user.Email).Limit(1).Scan(r.Context())
	if err == nil {
		utils.WriteResponse(w, http.StatusConflict, "Email already exists", nil, nil)
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to hash password")
		utils.WriteResponse(w, http.StatusInternalServerError, "Internal server error", nil, err)
		return
	}
	user.Password = string(hashedPassword)

	// Insert the new user
	_, err = db.DB.NewInsert().Model(&user).Exec(r.Context())
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to insert user")
		utils.WriteResponse(w, http.StatusInternalServerError, "Internal server error", nil, err)
		return
	}

	// Clear the password before sending the response
	user.Password = ""
	utils.WriteResponse(w, http.StatusCreated, "User created successfully", user, nil)
}

func Login(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Login called")
	if r.Method != http.MethodPost {
		utils.WriteResponse(w, http.StatusMethodNotAllowed, "Invalid request method", nil, nil)
		return
	}

	var loginDetails struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&loginDetails); err != nil {
		utils.WriteResponse(w, http.StatusBadRequest, "Invalid request payload", nil, err)
		return
	}

	if loginDetails.Email == "" {
		utils.WriteResponse(w, http.StatusBadRequest, "Email is required", nil, nil)
		return
	}
	if loginDetails.Password == "" {
		utils.WriteResponse(w, http.StatusBadRequest, "Password is required", nil, nil)
		return
	}

	var user models.User
	err := db.DB.NewSelect().Model(&user).Where("email = ?", loginDetails.Email).Limit(1).Scan(r.Context())
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to find user")
		utils.WriteResponse(w, http.StatusInternalServerError, "User not found", nil, err)
		return
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginDetails.Password)); err != nil {
		utils.WriteResponse(w, http.StatusUnauthorized, "Incorrect password", nil, err)
		return
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user.Name, user.ID)
	if err != nil {
		utils.WriteResponse(w, http.StatusInternalServerError, "Error generating token", nil, err)
		return
	}

	// Set the token as a cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   token,
		Expires: time.Now().Add(24 * time.Hour),
	})
	utils.WriteResponse(w, http.StatusOK, "Login successful", map[string]string{"token": token}, nil)
}

func GetUserWithAddresses(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteResponse(w, http.StatusMethodNotAllowed, "Invalid request method", nil, nil)
		return
	}
	// type Key string

	// const UserIDKey Key = "userID"
	// userID, ok := r.Context().Value(UserIDKey).(int64)
	// if !ok {
	// 	http.Error(w, "User ID not found or invalid type", http.StatusUnauthorized)
	// 	return
	// }

	var users []models.User

	err := db.DB.NewSelect().
		Model(&users).
		Relation("Addresses").
		Where("id = ?", 1). // Filter by user ID
		Scan(r.Context())
	if err != nil {
		utils.WriteResponse(w, http.StatusInternalServerError, "Failed to fetch users with addresses", nil, err)
		return
	}

	utils.WriteResponse(w, http.StatusOK, "Users and addresses fetched successfully", users, nil)
}
