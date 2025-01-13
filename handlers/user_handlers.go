package handlers

import (
	"filmatch/database"
	"filmatch/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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

// POST /user/login
func LoginUser(c *gin.Context) {
	// Let's get the email from the context object
	email, exists := c.Get("email")
	if !exists || email == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email not found in context"})
		return
	}

	// let's check if the email is a string
	emailStr, ok := email.(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
		return
	}

	// Let's verify if the user exists in the db
	var user model.User
	if err := database.DB.Where("email = ?", emailStr).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Let's create a new user if it doesn't exist
			newUser := model.User{
				Email: emailStr,
			}
			if err := database.DB.Create(&newUser).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user", "details": err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"message": "User created successfully",
				"user":    newUser,
			})
			return
		}

		// In case there's any other error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query user", "details": err.Error()})
		return
	}

	// If the user exists, let's return the user
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User logged in successfully",
		"user":    user,
	})
}
