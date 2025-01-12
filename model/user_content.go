package model

type UserContent struct {
	ID      uint    `gorm:"primaryKey" json:"id"`
	UserID  uint    `json:"user_id"`
	MovieID uint    `json:"movie_id"`
	Status  int     `json:"status"` // 1: Wants to watch, 2: Doesn't want to watch, 3: Seen, 4: Superlike
	User    User    `gorm:"foreignKey:UserID" json:"-"`
	Movie   Content `gorm:"foreignKey:MovieID" json:"-"`
}
