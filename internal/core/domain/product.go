package domain

import "gorm.io/gorm"

type Category struct {
	gorm.Model

	Name string `json:"name" gorm:"size:255;not null;unique"`
}

type Product struct {
	gorm.Model

	Name string `json:"name" gorm:"size:255;not null;unique"`
	Price float64 `json:"price" gorm:"not null"`
	Description string `json:"description" gorm:"size:255"`

	CategoryID uint `json:"category_id" gorm:"not null"`
	Category Category `json:"category" gorm:"foreignKey:CategoryID"`
}
