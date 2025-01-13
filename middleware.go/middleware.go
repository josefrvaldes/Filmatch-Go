package middleware

import (
	"context"
	"net/http"
	"strings"

	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
)

func FirebaseAuthMiddleware(client *auth.Client) gin.HandlerFunc {
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

		idToken := tokenParts[1]

		// Let's verify the token with Firebase
		token, err := client.VerifyIDToken(context.Background(), idToken)
		if err != nil {
			currentContext.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token", "details": err.Error()})
			currentContext.Abort()
			return
		}

		// Let's add the current user fields to the context
		currentContext.Set("uid", token.UID)
		currentContext.Set("email", token.Claims["email"])
		currentContext.Next()
	}
}
