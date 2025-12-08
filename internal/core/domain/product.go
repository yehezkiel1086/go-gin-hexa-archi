package domain

import "gorm.io/gorm"

type Category uint16

type Product struct {
	gorm.Model

	Name string `json:"name" gorm:"size:255;not null"`
	Description string `json:"description" gorm:"size:255;not null"`
	Price float64 `json:"price" gorm:"not null"`
}
