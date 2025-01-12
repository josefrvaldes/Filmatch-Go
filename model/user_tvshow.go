package model

type UserTVShow struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	UserID   uint   `json:"user_id"`
	TVShowID uint   `json:"tv_show_id"`
	Status   int    `json:"status"`
	User     User   `gorm:"foreignKey:UserID" json:"-"`
	TVShow   TVShow `gorm:"foreignKey:TVShowID" json:"-"`
}
