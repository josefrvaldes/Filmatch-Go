package model

type User struct {
	ID       uint    `gorm:"primaryKey" json:"id"`
	Email    string  `gorm:"unique;not null" json:"email"`
	Username *string `gorm:"unique;null" json:"username,omitempty"`
}
