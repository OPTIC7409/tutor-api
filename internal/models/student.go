package models

import (
	"gorm.io/gorm"
)

type Student struct {
	gorm.Model
	UserID   uint   `gorm:"not null"`
	User     User   `gorm:"foreignKey:UserID"`
	Age      int    `gorm:"not null"`
	Subjects string `gorm:"size:255;not null"`
	Location string `gorm:"size:255;not null"`
}
