package main

import (
	"context"
	"filmatch/database"
	"filmatch/firebase"
	"filmatch/handlers"
	"filmatch/interceptors"

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

	r.Use(interceptors.FirebaseAuthInterceptor(client))

	r.POST("/user/auth", handlers.PerformAuth)

	r.POST("/user/content", handlers.CreateUserVisit)

	r.GET("/user/:id/movie", handlers.GetUserVisitMoviesByStatus)
	r.GET("/user/:id/tv", handlers.GetUserVisitTVShowsByStatus)

	// Let's init the server
	r.Run(":9090")
}
