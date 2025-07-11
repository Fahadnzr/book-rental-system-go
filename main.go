package main

import (
	"book-rental-system/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.POST("/books", handlers.AddBook)
	r.GET("/books", handlers.ListBooks)
	// Register user handler
	r.POST("/users", handlers.AddUser)
	// Register rent and return handlers
	r.POST("/rent", handlers.RentBook)
	r.POST("/return", handlers.ReturnBook)

	r.Run(":8081")
}
