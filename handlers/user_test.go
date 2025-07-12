package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"book-rental-system/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// cleanupTestData resets the shared data stores for test isolation
func cleanupTestData() {
	booksMu.Lock()
	books = make(map[uuid.UUID]models.Book)
	booksMu.Unlock()

	usersMu.Lock()
	users = make(map[uuid.UUID]models.User)
	usersMu.Unlock()
}

func setupUserRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/users", AddUser)
	r.POST("/books", AddBook)
	return r
}

func TestAddUser_Valid(t *testing.T) {
	cleanupTestData() // Reset data before test
	r := setupUserRouter()
	body := map[string]string{"name": "Test User", "email": "test@example.com"}
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestAddUser_Invalid(t *testing.T) {
	cleanupTestData() // Reset data before test
	r := setupUserRouter()
	body := map[string]string{"name": ""}
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAddUser_InvalidEmail(t *testing.T) {
	cleanupTestData() // Reset data before test
	r := setupUserRouter()
	body := map[string]string{"name": "Test User", "email": "invalid-email"}
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAddUser_DuplicateEmail(t *testing.T) {
	cleanupTestData() // Reset data before test
	r := setupUserRouter()

	// First user
	body1 := map[string]string{"name": "Test User 1", "email": "test@example.com"}
	jsonBody1, _ := json.Marshal(body1)
	req1, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonBody1))
	req1.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusCreated, w1.Code)

	// Second user with same email
	body2 := map[string]string{"name": "Test User 2", "email": "test@example.com"}
	jsonBody2, _ := json.Marshal(body2)
	req2, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonBody2))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusConflict, w2.Code)
}
