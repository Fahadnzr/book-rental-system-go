package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Use usersMu, users, booksMu, books from other handler files
type RentInput struct {
	UserID string `json:"user_id" binding:"required,uuid"`
	BookID string `json:"book_id" binding:"required,uuid"`
}

func RentBook(c *gin.Context) {
	var input RentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	userID, err := uuid.Parse(input.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id"})
		return
	}
	bookID, err := uuid.Parse(input.BookID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book_id"})
		return
	}
	usersMu.RLock()
	_, userExists := users[userID]
	usersMu.RUnlock()
	if !userExists {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	booksMu.Lock()
	book, bookExists := books[bookID]
	if !bookExists {
		booksMu.Unlock()
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}
	if !book.Available {
		booksMu.Unlock()
		c.JSON(http.StatusConflict, gin.H{"error": "Book not available"})
		return
	}
	book.Available = false
	books[bookID] = book
	booksMu.Unlock()
	c.JSON(http.StatusOK, gin.H{"message": "Book rented successfully"})
}

func ReturnBook(c *gin.Context) {
	var input struct {
		BookID string `json:"book_id" binding:"required,uuid"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	bookID, err := uuid.Parse(input.BookID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book_id"})
		return
	}
	booksMu.Lock()
	book, bookExists := books[bookID]
	if !bookExists {
		booksMu.Unlock()
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}
	if book.Available {
		booksMu.Unlock()
		c.JSON(http.StatusConflict, gin.H{"error": "Book is not rented"})
		return
	}
	book.Available = true
	books[bookID] = book
	booksMu.Unlock()
	c.JSON(http.StatusOK, gin.H{"message": "Book returned successfully"})
}
