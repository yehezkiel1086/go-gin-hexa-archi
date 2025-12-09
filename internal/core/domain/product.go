package domain

import "gorm.io/gorm"

type Product struct {
	gorm.Model

	Name string `json:"name" gorm:"size:255;not null"`
	Description string `json:"description" gorm:"size:255;not null"`
	Price float64 `json:"price" gorm:"not null"`

	CategoryID uint `json:"category_id,omitempty" gorm:"not null"`
	Category Category `json:"category" gorm:"foreignKey:CategoryID"`
}
