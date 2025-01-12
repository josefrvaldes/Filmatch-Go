package model

import "gorm.io/gorm"

type Content struct {
	ID               uint    `gorm:"primaryKey" json:"id"`
	Adult            bool    `json:"adult"`
	BackdropPath     string  `json:"backdrop_path"`
	GenreIDs         []int   `gorm:"-" json:"genre_ids"` // It's not automatically saved in the db
	GenreIDsRaw      string  `json:"-"`                  // saved as string into the db
	OriginalLanguage string  `json:"original_language"`
	Overview         string  `json:"overview"`
	Popularity       float64 `json:"popularity"`
	PosterPath       string  `json:"poster_path"`
	VoteAverage      float64 `json:"vote_average"`
	VoteCount        int     `json:"vote_count"`

	// Specific fields for Movie
	OriginalTitle *string `json:"original_title,omitempty"`
	ReleaseDate   *string `json:"release_date,omitempty"`
	Title         *string `json:"title,omitempty"`
	Video         *bool   `json:"video,omitempty"`

	// Specific fields for TV Show
	OriginCountry []string `gorm:"-" json:"origin_country,omitempty"`
	OriginRaw     string   `json:"-"` // Guarda como string para SQL
	OriginalName  *string  `json:"original_name,omitempty"`
	FirstAirDate  *string  `json:"first_air_date,omitempty"`
	Name          *string  `json:"name,omitempty"`
}

// Before saving the content, we serialize the complex fields
func (c *Content) BeforeSave(tx *gorm.DB) (err error) {
	c.GenreIDsRaw = toJSON(c.GenreIDs)
	c.OriginRaw = toJSON(c.OriginCountry)
	return nil
}

// After finding the content, we deserialize the complex fields
func (c *Content) AfterFind(tx *gorm.DB) (err error) {
	c.GenreIDs = fromJSONIntSlice(c.GenreIDsRaw)
	c.OriginCountry = fromJSONStringSlice(c.OriginRaw)
	return nil
}
