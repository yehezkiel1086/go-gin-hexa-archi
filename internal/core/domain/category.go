package domain

import "gorm.io/gorm"

type Category struct {
	gorm.Model

	Name        string `gorm:"type:varchar(100);not null;unique" json:"name"`
	Description string `gorm:"type:text" json:"description"`
}
