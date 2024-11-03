package models

import (
	"gorm.io/gorm"
)

type Tutor struct {
	gorm.Model
	UserID          uint    `gorm:"not null"`
	User            User    `gorm:"foreignKey:UserID"`
	Subject         string  `gorm:"size:255;not null"`
	YearsExperience int     `gorm:"not null"`
	HourlyRate      float64 `gorm:"not null"`
	Location        string  `gorm:"size:255;not null"`
}
