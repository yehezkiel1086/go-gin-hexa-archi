package domain

import "gorm.io/gorm"

type Role uint16

const (
	AdminRole Role = 5150
	UserRole Role = 2001
)

type User struct {
	gorm.Model

	Name string `json:"name" gorm:"size:255;not null"`
	Email string `json:"email" gorm:"size:255;not null;unique"`
	Password string `json:"password" gorm:"size:255;not null"`
	Role Role `json:"role" gorm:"type:int;default:2001"`
}
