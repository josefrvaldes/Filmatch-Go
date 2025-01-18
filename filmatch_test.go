package main

import (
	"filmatch/handlers"
	"filmatch/interceptors"
	"filmatch/mocks"
	"filmatch/model"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&model.User{})
	db.Create(&model.User{
		Email: "test@example.com",
		UID:   "test-uid",
	})
	return db
}

func emptyTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&model.User{})
	return db
}

func TestHealthCheck(t *testing.T) {

	r := gin.Default()
	testDB := setupTestDB()

	SetupRoutes(r, testDB)

	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status %d but got %d", http.StatusOK, w.Code)
	}

	expectedBody := `{"status":"ok"}`
	if w.Body.String() != expectedBody {
		t.Fatalf("Expected body %s but got %s", expectedBody, w.Body.String())
	}
}

func TestFirebaseAuthInterceptor(t *testing.T) {

	mockAuthClient := &mocks.MockAuthClient{}

	databaseWithUser := setupTestDB()
	emptyDatabase := emptyTestDB()

	tests := []struct {
		name               string
		authHeader         string
		expectedStatusCode int
		expectedBody       string
		database           *gorm.DB
	}{
		{
			name:               "No Authorization Header",
			authHeader:         "",
			expectedStatusCode: http.StatusUnauthorized,
			expectedBody:       `{"error":"Authorization header missing"}`,
			database:           databaseWithUser,
		},
		{
			name:               "Invalid Token Format",
			authHeader:         "Bearer",
			expectedStatusCode: http.StatusUnauthorized,
			expectedBody:       `{"error":"Invalid Authorization header format"}`,
			database:           databaseWithUser,
		},
		{
			name:               "Invalid Token",
			authHeader:         "Bearer invalid-token",
			expectedStatusCode: http.StatusUnauthorized,
			expectedBody:       `{"details":"invalid token","error":"Invalid token"}`,
			database:           databaseWithUser,
		},
		{
			name:               "Valid Token, User Not Registered",
			authHeader:         "Bearer valid-token",
			expectedStatusCode: http.StatusUnauthorized,
			expectedBody:       `{"error":"User not registered in the system"}`,
			database:           emptyDatabase,
		},
		{
			name:               "Valid Token, User Registered",
			authHeader:         "Bearer valid-token",
			expectedStatusCode: http.StatusOK,
			expectedBody:       `{"message":"Access granted"}`,
			database:           databaseWithUser,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			testDB := tt.database

			// let's create a test route that uses the interceptor
			gin.SetMode(gin.TestMode)
			router := gin.Default()
			router.Use(interceptors.FirebaseAuthInterceptor(mockAuthClient, testDB))
			router.GET("/protected", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "Access granted"})
			})

			// let's create the http request
			req, _ := http.NewRequest("GET", "/protected", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			w := httptest.NewRecorder()

			// and send it
			router.ServeHTTP(w, req)

			// let's verify the statuses and bodies
			if w.Code != tt.expectedStatusCode {
				t.Errorf("Expected status %d but got %d", tt.expectedStatusCode, w.Code)
			}
			if strings.TrimSpace(w.Body.String()) != tt.expectedBody {
				t.Errorf("Expected body %s but got %s", tt.expectedBody, w.Body.String())
			}
		})
	}
}

func TestPerformAuth(t *testing.T) {

	databaseWithUser := setupTestDB()
	emptyDatabase := emptyTestDB()

	tests := []struct {
		name               string
		email              string
		uid                string
		expectedStatusCode int
		expectedBody       string
		database           *gorm.DB
	}{
		{
			name:               "User Not Found, Creates New User",
			email:              "newuser@example.com",
			uid:                "new-uid",
			expectedStatusCode: http.StatusOK,
			expectedBody:       `{"message":"User created successfully","success":true,"user":{"id":1,"email":"newuser@example.com","uid":"new-uid"}}`,
			database:           emptyDatabase,
		},
		{
			name:               "User Already Exists",
			email:              "test@example.com",
			uid:                "test-uid",
			expectedStatusCode: http.StatusOK,
			expectedBody:       `{"message":"User logged in successfully","success":true,"user":{"id":1,"email":"test@example.com","uid":"test-uid"}}`,
			database:           databaseWithUser,
		},
		{
			name:               "Missing Email in Context",
			email:              "",
			uid:                "test-uid",
			expectedStatusCode: http.StatusUnauthorized,
			expectedBody:       `{"error":"Email not found in context"}`,
			database:           emptyDatabase,
		},
		{
			name:               "Missing UID in Context",
			email:              "test@example.com",
			uid:                "",
			expectedStatusCode: http.StatusUnauthorized,
			expectedBody:       `{"error":"Uid not found in context"}`,
			database:           emptyDatabase,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDB := tt.database

			// let's create a test route that uses the handler
			gin.SetMode(gin.TestMode)
			router := gin.Default()
			router.POST("/user/auth", func(c *gin.Context) {
				// let's mock the context with the test data
				if tt.email != "" {
					c.Set("email", tt.email)
				}
				if tt.uid != "" {
					c.Set("uid", tt.uid)
				}
				handlers.PerformAuth(c, testDB)
			})

			// creating the request...
			req, _ := http.NewRequest("POST", "/user/auth", nil)
			w := httptest.NewRecorder()

			// and sending it
			router.ServeHTTP(w, req)

			// let's verify the statuses and bodies
			if w.Code != tt.expectedStatusCode {
				t.Errorf("Expected status %d but got %d", tt.expectedStatusCode, w.Code)
			}
			if strings.TrimSpace(w.Body.String()) != tt.expectedBody {
				t.Errorf("Expected body %s but got %s", tt.expectedBody, w.Body.String())
			}
		})
	}
}
