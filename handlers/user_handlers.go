package handlers

import (
	"filmatch/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// POST /user/auth
func PerformAuth(c *gin.Context, db *gorm.DB) {
	// Let's get the email from the context object
	email, exists := c.Get("email")
	if !exists || email == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email not found in context"})
		return
	}

	// Let's get the uid from the context object
	uid, exists := c.Get("uid")
	if !exists || uid == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Uid not found in context"})
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
	if err := db.Where("email = ?", emailStr).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Let's create a new user if it doesn't exist
			newUser := model.User{
				Email: emailStr,
				UID:   uid.(string),
			}
			if err := db.Create(&newUser).Error; err != nil {
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
