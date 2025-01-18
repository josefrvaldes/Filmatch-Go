package main

import (
	"context"
	"filmatch/database"
	"filmatch/firebase"
	"filmatch/handlers"
	"filmatch/interceptors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func main() {
	// Connect to the db
	database.ConnectDatabase()

	// let's init the firebase auth client
	firebase.InitFirebase()
	firebaseClient, err := firebase.App.Auth(context.Background())
	authClient := firebase.NewFirebaseAuthClient(firebaseClient)

	if err != nil {
		panic("Failed to connect to Firebase")
	}

	// Let's setup routes
	engine := gin.Default()

	// adding interceptors
	engine.Use(interceptors.VersionInterceptor())
	engine.Use(interceptors.FirebaseAuthInterceptor(authClient, database.DB))

	SetupRoutes(engine, database.DB)

	// Let's init the server
	engine.Run(":9090")
}

func SetupRoutes(engine *gin.Engine, db *gorm.DB) {
	engine.POST("/user/auth", func(c *gin.Context) {
		handlers.PerformAuth(c, db)
	})

	engine.POST("/user/content", func(c *gin.Context) {
		handlers.CreateUserVisit(c, db)
	})

	engine.GET("/user/:id/movie", func(c *gin.Context) {
		handlers.GetUserVisitMoviesByStatus(c, db)
	})

	engine.GET("/user/:id/tv", func(c *gin.Context) {
		handlers.GetUserVisitTVShowsByStatus(c, db)
	})

	engine.GET("/health", handlers.HealthCheck)
}
