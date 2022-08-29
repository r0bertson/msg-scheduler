package models

import (
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	Subject string `gorm:"size:200;not null;"`
	Content string `gorm:"size:5000;not null"`
}

type SentMessage struct {
	gorm.Model
	UserID    uint
	MessageID uint
}
