package domain

import (
	"time"

	"gorm.io/gorm"
)

type Role uint16

const (
	AdminRole Role = 5150
	UserRole Role = 2001
)

type User struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// credentials
	Username string `json:"username"`
	Password string `json:"password"`
	Role string `json:"role"`

	// details
	Fullname string `json:"fullname"`
	Email string `json:"email"`
}
