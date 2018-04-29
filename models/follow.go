package models

type Follow struct {
	BaseModel
	FollowBy int `gorm:"column:FollowBy"`
	FollowTo int `gorm:"column:FollowTo"`
}

func (Follow) TableName() string {
	return "Follows"
}
