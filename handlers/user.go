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
	users   = make(map[uuid.UUID]models.User)
	usersMu sync.RWMutex
)

type AddUserInput struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

func AddUser(c *gin.Context) {
	var input AddUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}
	// validation for empty name
	if len(strings.TrimSpace(input.Name)) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Name cannot be empty or just spaces"})
		return
	}
	// Length validation
	if len(strings.TrimSpace(input.Name)) > 100 || len(strings.TrimSpace(input.Email)) > 254 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Name must be <= 100 characters and Email must be <= 254 characters"})
		return
	}
	// Duplicate email check and user creation (atomic operation)
	emailToCheck := strings.ToLower(strings.TrimSpace(input.Email))
	usersMu.Lock()
	defer usersMu.Unlock()

	// Check for duplicates within the write lock
	for _, u := range users {
		if strings.ToLower(strings.TrimSpace(u.Email)) == emailToCheck {
			c.JSON(http.StatusConflict, gin.H{"message": "A user with the given email already exists. Please use a different email address."})
			return
		}
	}

	user := models.User{
		ID:    uuid.New(),
		Name:  input.Name,
		Email: input.Email,
	}
	users[user.ID] = user
	c.JSON(http.StatusCreated, user)
}
