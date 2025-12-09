package domain

import "gorm.io/gorm"

type Category struct {
	gorm.Model

	Name string `json:"name" gorm:"size:255;not null"`
}
