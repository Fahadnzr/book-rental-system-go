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

// cleanupBookTestData resets the shared data stores for test isolation
func cleanupBookTestData() {
	booksMu.Lock()
	books = make(map[uuid.UUID]models.Book)
	booksMu.Unlock()

	usersMu.Lock()
	users = make(map[uuid.UUID]models.User)
	usersMu.Unlock()
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/books", AddBook)
	r.GET("/books", ListBooks)
	return r
}

func TestAddBook_Valid(t *testing.T) {
	cleanupBookTestData() // Reset data before test
	r := setupRouter()
	body := map[string]string{"title": "Test Book", "author": "Author"}
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestAddBook_Invalid(t *testing.T) {
	cleanupBookTestData() // Reset data before test
	r := setupRouter()
	body := map[string]string{"title": ""}
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestListBooks(t *testing.T) {
	cleanupBookTestData() // Reset data before test
	r := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/books", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
