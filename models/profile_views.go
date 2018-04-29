package models

type ProfileView struct {
	BaseModel
	ViewBy int `gorm:"column:ViewBy"`
	ViewTo int `gorm:"column:ViewTo"`
}

func (ProfileView) TableName() string {
	return "ProfileViews"
}
