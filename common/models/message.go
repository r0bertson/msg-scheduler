package models

import (
	"gorm.io/gorm"
	"time"
)

type Message struct {
	gorm.Model
	Subject string `gorm:"size:200;not null;"`
	Content string `gorm:"size:5000;not null"`
}

func (msg *Message) SaveMessage(db *gorm.DB) (*Message, error) {
	if err := db.Debug().Create(&msg).Error; err != nil {
		return &Message{}, err
	}
	return msg, nil
}

func (msg *Message) UpdateMessage(db *gorm.DB) (*Message, error) {
	db = db.Debug().Model(&Message{}).Where("id = ?", msg.ID).Take(&Message{}).UpdateColumns(
		map[string]interface{}{
			"content":    msg.Content,
			"subject":    msg.Subject,
			"updated_at": time.Now(),
		},
	)
	if db.Error != nil {
		return &Message{}, db.Error
	}
	// This is the display the updated message
	if err := db.Debug().Model(&Message{}).Where("id = ?", msg.ID).Take(&msg).Error; err != nil {
		return &Message{}, err
	}
	return msg, nil
}
