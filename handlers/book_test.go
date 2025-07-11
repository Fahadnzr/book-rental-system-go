package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/books", AddBook)
	r.GET("/books", ListBooks)
	return r
}

func TestAddBook_Valid(t *testing.T) {
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
	r := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/books", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
