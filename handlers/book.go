package handlers

import (
	"net/http"
	"strings"
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
	//  Validation for empty or whitespace-only fields
	if len(strings.TrimSpace(input.Title)) == 0 || len(strings.TrimSpace(input.Author)) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Title and Author cannot be empty"})
		return
	}
	// Duplicate title check
	titleToCheck := strings.ToLower(strings.TrimSpace(input.Title))
	booksMu.RLock()
	for _, b := range books {
		if strings.ToLower(strings.TrimSpace(b.Title)) == titleToCheck {
			booksMu.RUnlock()
			c.JSON(http.StatusConflict, gin.H{"message": "A book with the given title already exists. Please use a different title."})
			return
		}
	}
	booksMu.RUnlock()
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
