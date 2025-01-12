package model

type TVShow struct {
	ID            uint              `gorm:"primaryKey" json:"-"`
	Content       `gorm:"embedded"` // Common fields inherited from Content
	OriginCountry []string          `gorm:"-" json:"origin_country,omitempty"`
	OriginRaw     string            `json:"-"` // Saved as string
	OriginalName  *string           `json:"original_name,omitempty"`
	FirstAirDate  *string           `json:"first_air_date,omitempty"`
	Name          *string           `json:"name,omitempty"`
}
