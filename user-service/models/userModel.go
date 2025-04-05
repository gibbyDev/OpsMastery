package models

import (
    "time"
	"gorm.io/gorm"
)

type User struct {
    gorm.Model
    Email             string `json:"email" gorm:"unique"`
    Password          string `json:"password"`
    Name              string `json:"name"`
    Role              string `json:"role" gorm:"default:Admin"`
    Active            bool   `json:"active" gorm:"default:false"`
    VerificationToken string `json:"-"`
    ResetToken        string `json:"-" "`
    ResetTokenExpiry  time.Time `json:"-"`
}
