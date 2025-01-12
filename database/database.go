package database

import (
	"filmatch/model"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
	// Cargar variables de entorno desde el archivo .env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found. Loading environment variables from the system")
	}
}

func ConnectDatabase() {
	// let's read the environment variables
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")

	if dbUser == "" || dbPass == "" || dbHost == "" || dbName == "" {
		log.Fatal("Database environment variables not set")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPass, dbHost, dbName)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to MySQL database:", err)
	}

	// Migrate Movie
	if err := DB.AutoMigrate(&model.Movie{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	// Migrate TvShow
	if err := DB.AutoMigrate(&model.TVShow{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Migrate UserContent
	if err := DB.AutoMigrate(&model.UserMovie{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Migrate UserContent
	if err := DB.AutoMigrate(&model.UserTVShow{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	// Migrate User
	if err := DB.AutoMigrate(&model.User{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
}
