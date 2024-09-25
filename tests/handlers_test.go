package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"gologin/internal/db"
	"gologin/internal/handlers"
	"gologin/internal/middleware"
	"gologin/internal/models"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	// Set up the database connection
	if err := db.SetupTestDB(); err != nil {
		panic("Failed to connect to test database: " + err.Error())
	}

	db.Migrate()

	// Run all the tests
	code := m.Run()

	// Disconnect and cleanup after tests
	db.Disconnect() // Close the DB connection when done

	// Exit with the result of m.Run()
	os.Exit(code)
}

func TestSignUp_ExistingEmail(t *testing.T) {
	// First create a user
	user := models.User{Name: "Test User", Email: "existing@example.com", Password: "password123"}
	body, _ := json.Marshal(user)
	req := httptest.NewRequest("POST", "/signup", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	handlers.SignUp(w, req)

	// Now try to create the same user again
	req = httptest.NewRequest("POST", "/signup", bytes.NewBuffer(body))
	w = httptest.NewRecorder()

	handlers.SignUp(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusConflict, res.StatusCode)
}

// func TestSignUp_Success(t *testing.T) {
// 	// Set up a test request
// 	user := models.User{Name: "Test User", Email: "test11@example11.com", Password: "password123"}
// 	body, _ := json.Marshal(user)
// 	req := httptest.NewRequest("POST", "/signup", bytes.NewBuffer(body))
// 	w := httptest.NewRecorder()

// 	handlers.SignUp(w, req)

// 	res := w.Result()
// 	assert.Equal(t, http.StatusCreated, res.StatusCode)
// }

// func TestSignUp_MissingFields(t *testing.T) {
// 	// Test missing name
// 	user := models.User{Email: "test11@example11.com", Password: "password123"}
// 	body, _ := json.Marshal(user)
// 	req := httptest.NewRequest("POST", "/signup", bytes.NewBuffer(body))
// 	w := httptest.NewRecorder()

// 	handlers.SignUp(w, req)

// 	res := w.Result()
// 	assert.Equal(t, http.StatusBadRequest, res.StatusCode)

// }

// func TestLogin_Success(t *testing.T) {
// 	// Set up a login test case
// 	loginDetails := `{"email":"test11@example11.com","password":"password123"}`
// 	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer([]byte(loginDetails)))
// 	w := httptest.NewRecorder()

// 	handlers.Login(w, req)

// 	res := w.Result()
// 	fmt.Println(res)
// 	assert.Equal(t, http.StatusOK, res.StatusCode)

// }

func TestLogin_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer([]byte("invalid json")))
	w := httptest.NewRecorder()

	handlers.Login(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestGetUserProfile_Unauthorized(t *testing.T) {
	req := httptest.NewRequest("GET", "/profile", nil)
	w := httptest.NewRecorder()

	handlers.GetUsersWithAddresses(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
}

func TestGetUserProfile_NonExistentUser(t *testing.T) {
	// Set up a request with a context missing a user ID
	req := httptest.NewRequest("GET", "/profile", nil)
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, int64(99999)) // Assuming 99999 does not exist
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handlers.GetUsersWithAddresses(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode) // Update based on your logic
}
