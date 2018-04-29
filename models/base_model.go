package models

import "time"

type BaseModel struct {
	ID        uint       `gorm:"column:ID;primary_key"`
	CreatedAt time.Time  `gorm:"column:CreatedAt"`
	UpdatedAt time.Time  `gorm:"column:UpdatedAt"`
	DeletedAt *time.Time `gorm:"column:DeletedAt;index"`
}
