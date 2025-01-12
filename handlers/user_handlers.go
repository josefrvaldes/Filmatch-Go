package handlers

import (
	"filmatch/database"
	"filmatch/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

// POST /user
func CreateUser(context *gin.Context) {
	var user model.User
	if err := context.ShouldBindJSON(&user); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Let's save into the db
	if err := database.DB.Create(&user).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save user"})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"data": user})
}
