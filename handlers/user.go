package handlers

import (
	"net/http"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	user := models.User{
		ID:    uuid.New(),
		Name:  input.Name,
		Email: input.Email,
	}
	usersMu.Lock()
	users[user.ID] = user
	usersMu.Unlock()
	c.JSON(http.StatusCreated, user)
}
