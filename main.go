package main

import (
	"filmatch/database"
	"filmatch/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	// Connect to the db
	database.ConnectDatabase()

	// Let's setup routes
	r := gin.Default()
	r.POST("/movies", handlers.CreateContent)

	r.POST("/user", handlers.CreateUser)

	// Let's init the server
	r.Run(":8080")
}
