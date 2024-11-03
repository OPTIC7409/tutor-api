package models

import (
	"gorm.io/gorm"
)

type Chat struct {
	gorm.Model
	Participants []User `gorm:"many2many:chat_participants;"`
	Messages     []Message
}

type Message struct {
	gorm.Model
	ChatID   uint
	SenderID uint
	Sender   User
	Content  string
}
