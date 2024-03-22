package models

type Comment struct {
	GormModel
	UserID  uint   `gorm:"not null" json:"user_id"`
	PhotoID uint   `gorm:"not null" json:"photo_id" form:"photo_id"`
	Message string `gorm:"not null" json:"message" form:"message" valid:"required~Your message is required"`
	User    User   `gorm:"foreignKey:UserID" json:"User"`
	Photo   Photo  `gorm:"foreignKey:PhotoID" json:"Photo"`
}