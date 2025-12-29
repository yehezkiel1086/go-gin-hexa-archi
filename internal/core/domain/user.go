package domain

import "gorm.io/gorm"

type Role uint16

const (
	UserRole Role = 2001
	AdminRole Role = 5150
)

type User struct {
	gorm.Model

	Email string `json:"email" gorm:"size:255;unique;not null"`
	Password string `json:"password" gorm:"size:255;not null"`
	Role Role `json:"role" gorm:"not null;default:2001"`
	Name string `json:"name" gorm:"size:255;not null"`
}
