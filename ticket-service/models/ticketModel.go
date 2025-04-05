package models

import (
    "gorm.io/gorm"
)

type Ticket struct {
    gorm.Model
    Title       string `json:"title" gorm:"not null"`
    Description string `json:"description"`
    Priority    string `json:"priority" gorm:"not null;default:'Medium'"` // Priority field
    UserID      uint   `json:"user_id"`                                   // Foreign key for the creator
    // Creator     User   `json:"creator" gorm:"foreignKey:UserID"`          // Relationship with the User model
    ClientID    uint   `json:"client_id"`
}