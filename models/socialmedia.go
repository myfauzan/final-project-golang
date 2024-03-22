package models

type SocialMedia struct {
	GormModel
	Name           string `gorm:"not null;type:varchar(100)" json:"name" form:"name" valid:"required~Your name is required"`
	SocialMediaUrl string `gorm:"not null" json:"social_media_url" form:"social_media_url" valid:"required~Your social_media_url is required"`
	UserID         uint   `gorm:"not null" json:"user_id"`
	User           User   `gorm:"foreignKey:UserID" json:"User"`
}