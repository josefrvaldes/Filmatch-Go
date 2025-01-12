package handlers

import (
	"net/http"

	"filmatch/database"
	"filmatch/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateUserContent(c *gin.Context) {
	var input struct {
		User   model.User    `json:"user"`
		Movie  *model.Movie  `json:"movie,omitempty"`
		TVShow *model.TVShow `json:"tv_show,omitempty"`
		Status int           `json:"status"`
	}

	// Parse JSON
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify or create the user
	var user model.User
	if err := database.DB.Where("email = ?", input.User.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find user"})
		return
	}

	// Verifiy content type (movie or tv_show)
	if input.Movie != nil {
		processMovie(c, user, input.Movie, input.Status)
	} else if input.TVShow != nil {
		processTVShow(c, user, input.TVShow, input.Status)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No valid content provided"})
	}
}

func processMovie(c *gin.Context, user model.User, movie *model.Movie, status int) {
	var existingMovie model.Movie
	if err := database.DB.Where("tmdb_id = ?", movie.TMDBID).First(&existingMovie).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Serialize GenreIDs
			movie.GenreIDsRaw = model.ToJSON(movie.GenreIDs)
			if err := database.DB.Create(movie).Error; err != nil {
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
	if err := database.DB.Where("user_id = ? AND movie_id = ?", user.ID, movie.ID).First(&userMovie).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			userMovie = model.UserMovie{
				UserID:  user.ID,
				MovieID: movie.ID,
				Status:  status,
			}
			if err := database.DB.Create(&userMovie).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user-movie relation"})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find user-movie relation"})
			return
		}
	} else {
		userMovie.Status = status
		if err := database.DB.Save(&userMovie).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user-movie relation"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "User-movie relation processed successfully"})
}

func processTVShow(c *gin.Context, user model.User, tvShow *model.TVShow, status int) {
	var existingTVShow model.TVShow
	if err := database.DB.Where("tmdb_id = ?", tvShow.TMDBID).First(&existingTVShow).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Serialize GenreIDs and OriginCountry
			tvShow.GenreIDsRaw = model.ToJSON(tvShow.GenreIDs)
			tvShow.OriginRaw = model.ToJSON(tvShow.OriginCountry)
			if err := database.DB.Create(tvShow).Error; err != nil {
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
	if err := database.DB.Where("user_id = ? AND tv_show_id = ?", user.ID, tvShow.ID).First(&userTVShow).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			userTVShow = model.UserTVShow{
				UserID:   user.ID,
				TVShowID: tvShow.ID,
				Status:   status,
			}
			if err := database.DB.Create(&userTVShow).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user-tv_show relation"})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find user-tv_show relation"})
			return
		}
	} else {
		userTVShow.Status = status
		if err := database.DB.Save(&userTVShow).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user-tv_show relation"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "User-tv_show relation processed successfully"})
}

func GetUserMoviesByStatus(c *gin.Context) {
	// Let's get the id from the path
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	// Parse JSON to get the status
	var input struct {
		Status int `json:"status"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Let's make the query
	var movies []model.Movie
	err := database.DB.Joins("JOIN user_movies ON user_movies.movie_id = movies.id").
		Where("user_movies.user_id = ? AND user_movies.status = ?", userID, input.Status).
		Find(&movies).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch movies"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"movies": movies})
}

func GetUserTVShowsByStatus(c *gin.Context) {
	// Let's get the id from the path
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	// Parse JSON to get the status
	var input struct {
		Status int `json:"status"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Let's make the query
	var tvShows []model.TVShow
	err := database.DB.Joins("JOIN user_tv_shows ON user_tv_shows.tv_show_id = tv_shows.id").
		Where("user_tv_shows.user_id = ? AND user_tv_shows.status = ?", userID, input.Status).
		Find(&tvShows).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch TV shows"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tv_shows": tvShows})
}