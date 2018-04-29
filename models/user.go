package models

type User struct {
	BaseModel
	Username string `gorm:"column:Username"`
	Password string `gorm:"column:Password"`
	Email    string `gorm:"column:Email"`
	Bio      string `gorm:"column:Bio"`
}
