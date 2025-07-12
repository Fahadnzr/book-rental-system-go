package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RentInput struct {
	UserID string `json:"user_id" binding:"required,uuid"`
	BookID string `json:"book_id" binding:"required,uuid"`
}

func RentBook(c *gin.Context) {
	var input RentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}
	userID, err := uuid.Parse(input.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user_id"})
		return
	}
	bookID, err := uuid.Parse(input.BookID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid book_id"})
		return
	}
	usersMu.RLock()
	_, userExists := users[userID]
	usersMu.RUnlock()
	if !userExists {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}
	booksMu.Lock()
	defer booksMu.Unlock()

	book, bookExists := books[bookID]
	if !bookExists {
		c.JSON(http.StatusNotFound, gin.H{"message": "Book not found"})
		return
	}
	if !book.Available {
		c.JSON(http.StatusConflict, gin.H{"message": "Book not available"})
		return
	}
	book.Available = false
	books[bookID] = book
	c.JSON(http.StatusOK, gin.H{"message": "Book rented successfully"})
}

func ReturnBook(c *gin.Context) {
	var input struct {
		UserID string `json:"user_id" binding:"required,uuid"`
		BookID string `json:"book_id" binding:"required,uuid"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}
	userID, err := uuid.Parse(input.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user_id"})
		return
	}
	bookID, err := uuid.Parse(input.BookID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid book_id"})
		return
	}
	// Check if user exists
	usersMu.RLock()
	_, userExists := users[userID]
	usersMu.RUnlock()
	if !userExists {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}
	booksMu.Lock()
	defer booksMu.Unlock()

	book, bookExists := books[bookID]
	if !bookExists {
		c.JSON(http.StatusNotFound, gin.H{"message": "Book not found"})
		return
	}
	if book.Available {
		c.JSON(http.StatusConflict, gin.H{"message": "Book is not rented"})
		return
	}
	book.Available = true
	books[bookID] = book
	c.JSON(http.StatusOK, gin.H{"message": "Book returned successfully"})
}
