package model

type Movie struct {
	ID            uint              `gorm:"primaryKey" json:"-"`
	Content       `gorm:"embedded"` // Common fields inherited from Content
	OriginalTitle *string           `json:"original_title,omitempty"`
	ReleaseDate   *string           `json:"release_date,omitempty"`
	Title         *string           `json:"title,omitempty"`
	Video         *bool             `json:"video,omitempty"`
}
