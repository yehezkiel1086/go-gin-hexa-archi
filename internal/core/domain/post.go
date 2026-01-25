package domain

import "gorm.io/gorm"

type Post struct {
	gorm.Model

	CategoryID uint     `gorm:"not null" json:"category_id"`
	Category   Category `gorm:"foreignKey:CategoryID" json:"category"`

	Title     string `gorm:"type:varchar(255);not null" json:"title"`
	Content   string `gorm:"type:text;not null" json:"content"`
	Published bool   `gorm:"default:false" json:"published"`

	UserID uint `gorm:"not null" json:"user_id"`
	User   User `gorm:"foreignKey:UserID" json:"user"`

	Slug string `gorm:"type:varchar(255);not null;unique" json:"slug"`
}
