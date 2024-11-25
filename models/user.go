package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string `json:"email" gorm:"unique;not null"`
	Password string 
	Name     string `json:"name"`
	Role     string `json:"role"`
}

