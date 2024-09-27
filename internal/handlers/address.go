package handlers

import (
	"encoding/json"
	"fmt"
	"gologin/internal/db"
	"gologin/internal/logger"
	"gologin/internal/models"
	"gologin/internal/utils"
	"net/http"
	"strconv"
	"strings"
)

func AddAddress(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteResponse(w, http.StatusMethodNotAllowed, "Invalid request method", nil, nil)
		return
	}

	var address models.Address
	if err := json.NewDecoder(r.Body).Decode(&address); err != nil {
		utils.WriteResponse(w, http.StatusBadRequest, "Invalid request payload", nil, err)
		return
	}

	// Check for required fields
	if address.UserID == 0 {
		utils.WriteResponse(w, http.StatusBadRequest, "User ID is required", nil, nil)
		return
	}
	if address.Street == "" {
		utils.WriteResponse(w, http.StatusBadRequest, "Street is required", nil, nil)
		return
	}
	if address.City == "" {
		utils.WriteResponse(w, http.StatusBadRequest, "City is required", nil, nil)
		return
	}
	if address.State == "" {
		utils.WriteResponse(w, http.StatusBadRequest, "State is required", nil, nil)
		return
	}
	if address.ZipCode == "" {
		utils.WriteResponse(w, http.StatusBadRequest, "Zip code is required", nil, nil)
		return
	}
	if address.Country == "" {
		utils.WriteResponse(w, http.StatusBadRequest, "Country is required", nil, nil)
		return
	}

	// Insert the new address
	_, err := db.DB.NewInsert().Model(&address).Exec(r.Context())
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to insert address")
		utils.WriteResponse(w, http.StatusInternalServerError, "Internal server error", nil, err)
		return
	}

	utils.WriteResponse(w, http.StatusCreated, "Address added successfully", address, nil)
}

func UpdateAddress(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utils.WriteResponse(w, http.StatusMethodNotAllowed, "Invalid request method", nil, nil)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/users/address/")
	segments := strings.Split(path, "/")
	fmt.Println("segments", len(segments) == 0, segments[4])
	// Assuming the ID is the first segment after the prefix
	if len(segments) == 0 || segments[len(segments)-1] == "" {
		http.Error(w, "Address ID is required", http.StatusBadRequest)
		return
	}

	addressID, err := strconv.Atoi(segments[len(segments)-1]) // Convert ID from string to int
	if err != nil {
		utils.WriteResponse(w, http.StatusBadRequest, "Invalid address ID", nil, nil)
		return
	}
	fmt.Println("addressId", addressID)

	var address models.Address
	if err := json.NewDecoder(r.Body).Decode(&address); err != nil {
		utils.WriteResponse(w, http.StatusBadRequest, "Invalid request payload", nil, err)
		return
	}

	if address.UserID == 0 {
		utils.WriteResponse(w, http.StatusBadRequest, "User ID is required", nil, nil)
		return
	}
	if address.Street == "" {
		utils.WriteResponse(w, http.StatusBadRequest, "Street is required", nil, nil)
		return
	}
	if address.City == "" {
		utils.WriteResponse(w, http.StatusBadRequest, "City is required", nil, nil)
		return
	}
	if address.State == "" {
		utils.WriteResponse(w, http.StatusBadRequest, "State is required", nil, nil)
		return
	}
	if address.ZipCode == "" {
		utils.WriteResponse(w, http.StatusBadRequest, "Zip code is required", nil, nil)
		return
	}
	if address.Country == "" {
		utils.WriteResponse(w, http.StatusBadRequest, "Country is required", nil, nil)
		return
	}

	// Update the address in the database
	_, err = db.DB.NewUpdate().Model(&address).Where("id = ?", addressID).Exec(r.Context())
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to update address")
		utils.WriteResponse(w, http.StatusInternalServerError, "Internal server error", nil, err)
		return
	}

	utils.WriteResponse(w, http.StatusOK, "Address updated successfully", address, nil)
}
