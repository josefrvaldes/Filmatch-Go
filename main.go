package main

import (
	"context"
	"filmatch/database"
	"filmatch/firebase"
	"filmatch/handlers"
	"filmatch/middleware.go"

	"github.com/gin-gonic/gin"
)

func main() {
	// Connect to the db
	database.ConnectDatabase()

	firebase.InitFirebase()
	client, err := firebase.App.Auth(context.Background())

	if err != nil {
		panic("Failed to connect to Firebase")
	}

	// Let's setup routes
	r := gin.Default()

	r.Use(middleware.FirebaseAuthMiddleware(client))

	r.POST("/user/login", handlers.LoginUser)

	r.POST("/user/content", handlers.CreateUserContent)

	r.POST("/user", handlers.CreateUser)

	r.GET("/user/:id/movie", handlers.GetUserMoviesByStatus)
	r.GET("/user/:id/tv", handlers.GetUserTVShowsByStatus)

	// Let's init the server
	r.Run(":9090")
}
