package models

type Like struct {
	BaseModel
	PostID int `gorm:"column:PostID"`
	LikeBy int `gorm:"column:LikeBy"`
}
