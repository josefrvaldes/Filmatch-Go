package model

type UserMovie struct {
	ID      uint  `gorm:"primaryKey" json:"id"`
	UserID  uint  `json:"user_id"`
	MovieID uint  `json:"movie_id"`
	Status  int   `json:"status"`
	User    User  `gorm:"foreignKey:UserID" json:"-"`
	Movie   Movie `gorm:"foreignKey:MovieID" json:"-"`
}
