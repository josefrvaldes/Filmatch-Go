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
	r.POST("/user/content", handlers.CreateUserContent)

	r.POST("/user", handlers.CreateUser)

	r.GET("/user/:id/movie", handlers.GetUserMoviesByStatus)
	r.GET("/user/:id/tv", handlers.GetUserTVShowsByStatus)

	// Let's init the server
	r.Run(":9090")
}
