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

// cleanupRentTestData resets the shared data stores for test isolation
func cleanupRentTestData() {
	booksMu.Lock()
	books = make(map[uuid.UUID]models.Book)
	booksMu.Unlock()

	usersMu.Lock()
	users = make(map[uuid.UUID]models.User)
	usersMu.Unlock()
}

func setupRentRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/users", AddUser)
	r.POST("/books", AddBook)
	r.POST("/rent", RentBook)
	r.POST("/return", ReturnBook)
	return r
}

func TestRentBook_Valid(t *testing.T) {
	cleanupRentTestData() // Reset data before test
	r := setupRentRouter()

	// First create a user and book
	userBody := map[string]string{"name": "Test User", "email": "test@example.com"}
	userJson, _ := json.Marshal(userBody)
	userReq, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(userJson))
	userReq.Header.Set("Content-Type", "application/json")
	userW := httptest.NewRecorder()
	r.ServeHTTP(userW, userReq)

	bookBody := map[string]string{"title": "Test Book", "author": "Test Author"}
	bookJson, _ := json.Marshal(bookBody)
	bookReq, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(bookJson))
	bookReq.Header.Set("Content-Type", "application/json")
	bookW := httptest.NewRecorder()
	r.ServeHTTP(bookW, bookReq)

	// Extract user and book IDs from responses
	var userResponse map[string]interface{}
	json.Unmarshal(userW.Body.Bytes(), &userResponse)
	var bookResponse map[string]interface{}
	json.Unmarshal(bookW.Body.Bytes(), &bookResponse)

	// Now rent the book
	rentBody := map[string]string{
		"user_id": userResponse["id"].(string),
		"book_id": bookResponse["id"].(string),
	}
	rentJson, _ := json.Marshal(rentBody)
	rentReq, _ := http.NewRequest("POST", "/rent", bytes.NewBuffer(rentJson))
	rentReq.Header.Set("Content-Type", "application/json")
	rentW := httptest.NewRecorder()
	r.ServeHTTP(rentW, rentReq)

	assert.Equal(t, http.StatusOK, rentW.Code)
}

func TestRentBook_InvalidUUID(t *testing.T) {
	cleanupRentTestData() // Reset data before test
	r := setupRentRouter()
	body := map[string]string{"user_id": "invalid-uuid", "book_id": "invalid-uuid"}
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/rent", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRentBook_UserNotFound(t *testing.T) {
	cleanupRentTestData() // Reset data before test
	r := setupRentRouter()

	// Create a book first
	bookBody := map[string]string{"title": "Test Book", "author": "Test Author"}
	bookJson, _ := json.Marshal(bookBody)
	bookReq, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(bookJson))
	bookReq.Header.Set("Content-Type", "application/json")
	bookW := httptest.NewRecorder()
	r.ServeHTTP(bookW, bookReq)

	var bookResponse map[string]interface{}
	json.Unmarshal(bookW.Body.Bytes(), &bookResponse)

	// Try to rent with non-existent user
	rentBody := map[string]string{
		"user_id": "550e8400-e29b-41d4-a716-446655440000", // Non-existent UUID
		"book_id": bookResponse["id"].(string),
	}
	rentJson, _ := json.Marshal(rentBody)
	rentReq, _ := http.NewRequest("POST", "/rent", bytes.NewBuffer(rentJson))
	rentReq.Header.Set("Content-Type", "application/json")
	rentW := httptest.NewRecorder()
	r.ServeHTTP(rentW, rentReq)

	assert.Equal(t, http.StatusNotFound, rentW.Code)
}

func TestRentBook_BookNotFound(t *testing.T) {
	cleanupRentTestData() // Reset data before test
	r := setupRentRouter()

	// Create a user first
	userBody := map[string]string{"name": "Test User", "email": "test@example.com"}
	userJson, _ := json.Marshal(userBody)
	userReq, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(userJson))
	userReq.Header.Set("Content-Type", "application/json")
	userW := httptest.NewRecorder()
	r.ServeHTTP(userW, userReq)

	var userResponse map[string]interface{}
	json.Unmarshal(userW.Body.Bytes(), &userResponse)

	// Try to rent non-existent book
	rentBody := map[string]string{
		"user_id": userResponse["id"].(string),
		"book_id": "550e8400-e29b-41d4-a716-446655440000", // Non-existent UUID
	}
	rentJson, _ := json.Marshal(rentBody)
	rentReq, _ := http.NewRequest("POST", "/rent", bytes.NewBuffer(rentJson))
	rentReq.Header.Set("Content-Type", "application/json")
	rentW := httptest.NewRecorder()
	r.ServeHTTP(rentW, rentReq)

	assert.Equal(t, http.StatusNotFound, rentW.Code)
}

func TestReturnBook_Valid(t *testing.T) {
	cleanupRentTestData() // Reset data before test
	r := setupRentRouter()

	// First create a user and book
	userBody := map[string]string{"name": "Test User", "email": "test@example.com"}
	userJson, _ := json.Marshal(userBody)
	userReq, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(userJson))
	userReq.Header.Set("Content-Type", "application/json")
	userW := httptest.NewRecorder()
	r.ServeHTTP(userW, userReq)

	bookBody := map[string]string{"title": "Test Book", "author": "Test Author"}
	bookJson, _ := json.Marshal(bookBody)
	bookReq, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(bookJson))
	bookReq.Header.Set("Content-Type", "application/json")
	bookW := httptest.NewRecorder()
	r.ServeHTTP(bookW, bookReq)

	// Extract user and book IDs from responses
	var userResponse map[string]interface{}
	json.Unmarshal(userW.Body.Bytes(), &userResponse)
	var bookResponse map[string]interface{}
	json.Unmarshal(bookW.Body.Bytes(), &bookResponse)

	// Rent the book first
	rentBody := map[string]string{
		"user_id": userResponse["id"].(string),
		"book_id": bookResponse["id"].(string),
	}
	rentJson, _ := json.Marshal(rentBody)
	rentReq, _ := http.NewRequest("POST", "/rent", bytes.NewBuffer(rentJson))
	rentReq.Header.Set("Content-Type", "application/json")
	rentW := httptest.NewRecorder()
	r.ServeHTTP(rentW, rentReq)

	// Now return the book
	returnBody := map[string]string{
		"user_id": userResponse["id"].(string),
		"book_id": bookResponse["id"].(string),
	}
	returnJson, _ := json.Marshal(returnBody)
	returnReq, _ := http.NewRequest("POST", "/return", bytes.NewBuffer(returnJson))
	returnReq.Header.Set("Content-Type", "application/json")
	returnW := httptest.NewRecorder()
	r.ServeHTTP(returnW, returnReq)

	assert.Equal(t, http.StatusOK, returnW.Code)
}

func TestReturnBook_BookNotFound(t *testing.T) {
	cleanupRentTestData() // Reset data before test
	r := setupRentRouter()

	// Create a user first
	userBody := map[string]string{"name": "Test User", "email": "test@example.com"}
	userJson, _ := json.Marshal(userBody)
	userReq, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(userJson))
	userReq.Header.Set("Content-Type", "application/json")
	userW := httptest.NewRecorder()
	r.ServeHTTP(userW, userReq)

	var userResponse map[string]interface{}
	json.Unmarshal(userW.Body.Bytes(), &userResponse)

	// Try to return non-existent book
	returnBody := map[string]string{
		"user_id": userResponse["id"].(string),
		"book_id": "550e8400-e29b-41d4-a716-446655440000", // Non-existent UUID
	}
	returnJson, _ := json.Marshal(returnBody)
	returnReq, _ := http.NewRequest("POST", "/return", bytes.NewBuffer(returnJson))
	returnReq.Header.Set("Content-Type", "application/json")
	returnW := httptest.NewRecorder()
	r.ServeHTTP(returnW, returnReq)

	assert.Equal(t, http.StatusNotFound, returnW.Code)
}

func TestReturnBook_BookNotRented(t *testing.T) {
	cleanupRentTestData() // Reset data before test
	r := setupRentRouter()

	// Create a user and book
	userBody := map[string]string{"name": "Test User", "email": "test@example.com"}
	userJson, _ := json.Marshal(userBody)
	userReq, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(userJson))
	userReq.Header.Set("Content-Type", "application/json")
	userW := httptest.NewRecorder()
	r.ServeHTTP(userW, userReq)

	bookBody := map[string]string{"title": "Test Book", "author": "Test Author"}
	bookJson, _ := json.Marshal(bookBody)
	bookReq, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(bookJson))
	bookReq.Header.Set("Content-Type", "application/json")
	bookW := httptest.NewRecorder()
	r.ServeHTTP(bookW, bookReq)

	var userResponse map[string]interface{}
	json.Unmarshal(userW.Body.Bytes(), &userResponse)
	var bookResponse map[string]interface{}
	json.Unmarshal(bookW.Body.Bytes(), &bookResponse)

	// Try to return a book that's not rented
	returnBody := map[string]string{
		"user_id": userResponse["id"].(string),
		"book_id": bookResponse["id"].(string),
	}
	returnJson, _ := json.Marshal(returnBody)
	returnReq, _ := http.NewRequest("POST", "/return", bytes.NewBuffer(returnJson))
	returnReq.Header.Set("Content-Type", "application/json")
	returnW := httptest.NewRecorder()
	r.ServeHTTP(returnW, returnReq)

	assert.Equal(t, http.StatusConflict, returnW.Code)
}
