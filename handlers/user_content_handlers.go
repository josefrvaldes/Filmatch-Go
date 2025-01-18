package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"filmatch/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateUserVisit(c *gin.Context, db *gorm.DB) {
	var input struct {
		Movie  *model.Movie  `json:"movie,omitempty"`
		TVShow *model.TVShow `json:"tv_show,omitempty"`
		Status int           `json:"status"`
	}

	// let's verify if the user id is in the context. That means that this endpoint has been
	// protected by the middleware and therefore we don't need to check if the user exists again
	userId, exists := c.Get("userId")
	if !exists || userId == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "userId not found in context"})
		return
	}

	// Ensure userId is of type uint
	userIdUint, ok := userId.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "userId is not a valid uint"})
		return
	}

	// Parse JSON
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verifiy content type (movie or tv_show)
	if input.Movie != nil {
		createUserVisitMovie(c, db, userIdUint, input.Movie, input.Status)
	} else if input.TVShow != nil {
		createUserVisitTVShow(c, db, userIdUint, input.TVShow, input.Status)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No valid content provided"})
	}
}

func createUserVisitMovie(c *gin.Context, db *gorm.DB, userId uint, movie *model.Movie, status int) {
	var existingMovie model.Movie
	if err := db.Where("tmdb_id = ?", movie.TMDBID).First(&existingMovie).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Serialize GenreIDs
			movie.GenreIDsRaw = model.ToJSON(movie.GenreIDs)
			if err := db.Create(movie).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create movie"})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find movie"})
			return
		}
	} else {
		*movie = existingMovie // If the movie already exists, update the model in memory
	}

	// Many to Many relationship
	var userMovie model.UserMovie
	if err := db.Where("user_id = ? AND movie_id = ?", userId, movie.ID).First(&userMovie).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			userMovie = model.UserMovie{
				UserID:  userId,
				MovieID: movie.ID,
				Status:  status,
			}
			if err := db.Create(&userMovie).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user-movie relation"})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find user-movie relation"})
			return
		}
	} else {
		userMovie.Status = status
		if err := db.Save(&userMovie).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user-movie relation"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "User-movie relation processed successfully"})
}

func createUserVisitTVShow(c *gin.Context, db *gorm.DB, userId uint, tvShow *model.TVShow, status int) {
	var existingTVShow model.TVShow
	if err := db.Where("tmdb_id = ?", tvShow.TMDBID).First(&existingTVShow).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Serialize GenreIDs and OriginCountry
			tvShow.GenreIDsRaw = model.ToJSON(tvShow.GenreIDs)
			tvShow.OriginRaw = model.ToJSON(tvShow.OriginCountry)
			if err := db.Create(tvShow).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create TV show"})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find TV show"})
			return
		}
	} else {
		*tvShow = existingTVShow // Let's update the model in memory
	}

	// Many to Many relationship
	var userTVShow model.UserTVShow
	if err := db.Where("user_id = ? AND tv_show_id = ?", userId, tvShow.ID).First(&userTVShow).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			userTVShow = model.UserTVShow{
				UserID:   userId,
				TVShowID: tvShow.ID,
				Status:   status,
			}
			if err := db.Create(&userTVShow).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user-tv_show relation"})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find user-tv_show relation"})
			return
		}
	} else {
		userTVShow.Status = status
		if err := db.Save(&userTVShow).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user-tv_show relation"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "User-tv_show relation processed successfully"})
}

func getUserVisitsByStatus[T any](c *gin.Context, db *gorm.DB, tableName string, joinTable string, column string) {
	// Get the user ID from the path
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	// Parse the status from the JSON body
	var input struct {
		Status int `json:"status"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	resultsPerPage := 20

	if page < 1 {
		page = 1
	}
	offset := (page - 1) * resultsPerPage

	// Fetch total count
	var totalResults int64
	if err := db.Table(tableName).
		Joins(fmt.Sprintf("JOIN %s ON %s.%s = %s.id", joinTable, joinTable, column, tableName)).
		Where(fmt.Sprintf("%s.user_id = ? AND %s.status = ?", joinTable, joinTable), userID, input.Status).
		Count(&totalResults).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count results"})
		return
	}

	// Fetch paginated results
	var results []T
	if err := db.Table(tableName).
		Joins(fmt.Sprintf("JOIN %s ON %s.%s = %s.id", joinTable, joinTable, column, tableName)).
		Where(fmt.Sprintf("%s.user_id = ? AND %s.status = ?", joinTable, joinTable), userID, input.Status).
		Limit(resultsPerPage).Offset(offset).
		Find(&results).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch results"})
		return
	}

	// Calculate total pages
	totalPages := int((totalResults + int64(resultsPerPage) - 1) / int64(resultsPerPage))

	// Response
	c.JSON(http.StatusOK, gin.H{
		"page":             page,
		"results":          results,
		"results_per_page": resultsPerPage,
		"total_pages":      totalPages,
		"total_results":    totalResults,
	})
}

func GetUserVisitMoviesByStatus(c *gin.Context, db *gorm.DB) {
	getUserVisitsByStatus[model.Movie](c, db, "movies", "user_movies", "movie_id")
}

func GetUserVisitTVShowsByStatus(c *gin.Context, db *gorm.DB) {
	getUserVisitsByStatus[model.TVShow](c, db, "tv_shows", "user_tv_shows", "tv_show_id")
}
