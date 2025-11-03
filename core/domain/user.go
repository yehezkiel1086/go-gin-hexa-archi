package domain

import "gorm.io/gorm"

type Role uint16

const (
	AdminRole Role = 5150
	EmployeeRole Role = 2001
)

type User struct {
	gorm.Model
	
	Email    string `json:"email" gorm:"unique;not null;type:varchar(255);size:255"`
	Password string `json:"password" gorm:"not null;type:varchar(255);size:255"`
	Role     Role `json:"role" gorm:"default:2001"`

	Name     string `json:"name" gorm:"not null;type:varchar(255);size:255"`
}
