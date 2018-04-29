package models

import (
	"github.com/jinzhu/gorm"
)

type Follow struct {
	gorm.Model
	FollowBy int
	FollowTo int
}
