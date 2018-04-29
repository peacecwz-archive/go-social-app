package models

import (
	"github.com/jinzhu/gorm"
)

type Like struct {
	gorm.Model
	PostID int
	LikeBy int
}
