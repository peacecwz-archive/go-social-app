package models

import (
	"github.com/jinzhu/gorm"
)

type ProfileView struct {
	gorm.Model
	ViewBy int
	ViewTo int
}
