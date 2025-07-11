package handlers

import (
	"net/http"
	"sync"

	"book-rental-system/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var (
	books   = make(map[uuid.UUID]models.Book)
	booksMu sync.RWMutex
)

type AddBookInput struct {
	Title  string `json:"title" binding:"required"`
	Author string `json:"author" binding:"required"`
}

func AddBook(c *gin.Context) {
	var input AddBookInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	book := models.Book{
		ID:        uuid.New(),
		Title:     input.Title,
		Author:    input.Author,
		Available: true,
	}
	booksMu.Lock()
	books[book.ID] = book
	booksMu.Unlock()
	c.JSON(http.StatusCreated, book)
}

func ListBooks(c *gin.Context) {
	booksMu.RLock()
	defer booksMu.RUnlock()
	result := make([]models.Book, 0, len(books))
	for _, b := range books {
		result = append(result, b)
	}
	c.JSON(http.StatusOK, result)
}
