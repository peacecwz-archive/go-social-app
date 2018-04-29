package models

type Post struct {
	BaseModel
	Title     string `gorm:"column:Title"`
	Content   string `gorm:"column:Content"`
	CreatedBy int    `gorm:"column:CreatedBy"`
}
