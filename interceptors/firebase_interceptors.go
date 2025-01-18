package interceptors

import (
	"context"
	"filmatch/database"
	"filmatch/model"
	"net/http"
	"strings"

	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func FirebaseAuthInterceptor(client *auth.Client) gin.HandlerFunc {
	return func(currentContext *gin.Context) {
		// Let's extract the Authorization header
		authHeader := currentContext.GetHeader("Authorization")
		if authHeader == "" {
			currentContext.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			currentContext.Abort()
			return
		}

		// Let's check the format of the header
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			currentContext.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			currentContext.Abort()
			return
		}

		// Let's verify the token with Firebase
		idToken := tokenParts[1]
		token, err := client.VerifyIDToken(context.Background(), idToken)
		if err != nil {
			currentContext.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token", "details": err.Error()})
			currentContext.Abort()
			return
		}

		// Let's get the user email from the token
		email, ok := token.Claims["email"].(string)
		if !ok || email == "" {
			currentContext.JSON(http.StatusBadRequest, gin.H{"error": "Email not found in token"})
			currentContext.Abort()
			return
		}

		// If the endpoint is "/user/login", let's omit the database query, the user may not exist yet
		if currentContext.Request.URL.Path == "/user/auth" {
			currentContext.Set("uid", token.UID)
			currentContext.Set("email", email)
			currentContext.Next()
			return
		}

		// If it's any other endpoint, let's verify if the user exists in the database
		var user model.User
		if err := database.DB.Where("email = ?", email).First(&user).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				currentContext.JSON(http.StatusUnauthorized, gin.H{"error": "User not registered in the system"})
				currentContext.Abort()
				return
			}
			currentContext.JSON(http.StatusInternalServerError, gin.H{"error": "Database query failed", "details": err.Error()})
			currentContext.Abort()
			return
		}

		// Let's add the current user fields to the context
		currentContext.Set("userId", user.ID)
		currentContext.Set("uid", token.UID)
		currentContext.Set("email", token.Claims["email"])
		currentContext.Next()
	}
}
