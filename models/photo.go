package models

type Photo struct {
	GormModel
	Title    string `gorm:"not null" json:"title" form:"title" valid:"required~Your title is required"`
	Caption  string `gorm:"not null" json:"caption" form:"caption" valid:"required~Your caption is required"`
	PhotoUrl string `gorm:"not null" json:"photo_url" form:"photo_url" valid:"required~Your photo_url is required"`
	UserID   uint   `gorm:"not null" json:"user_id"`
	User     User   `gorm:"foreignKey:UserID" json:"User"`
}