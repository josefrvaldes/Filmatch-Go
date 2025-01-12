package handlers

import (
	"filmatch/database"
	"filmatch/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

// POST /content
func CreateContent(c *gin.Context) {
	var content model.Content
	if err := c.ShouldBindJSON(&content); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Let's save into the db
	if err := database.DB.Create(&content).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save content"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": content})
}
