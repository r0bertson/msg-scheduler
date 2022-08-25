package models

import (
	"errors"
	"gorm.io/gorm"
	"time"
)

type Message struct {
	gorm.Model
	Subject string `gorm:"size:200;not null;" json:"subject"`
	Content string `gorm:"size:5000;not null" json:"content"`
}

func (msg *Message) SaveMessage(db *gorm.DB) (*Message, error) {
	if err := db.Debug().Create(&msg).Error; err != nil {
		return &Message{}, err
	}
	return msg, nil
}

func (msg *Message) FindAllMessages(db *gorm.DB) (*[]Message, error) {
	messages := []Message{}
	if err := db.Debug().Model(&Message{}).Limit(100).Find(&messages).Error; err != nil {
		return &[]Message{}, err
	}
	return &messages, nil
}

func (msg *Message) FindMessageByID(db *gorm.DB, uid uint32) (*Message, error) {
	if err := db.Debug().Model(Message{}).Where("id = ?", uid).Take(&msg).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &Message{}, errors.New("message not found")
		}
		return &Message{}, err
	}

	return msg, nil
}

func (msg *Message) UpdateMessage(db *gorm.DB, uid uint32) (*Message, error) {
	db = db.Debug().Model(&Message{}).Where("id = ?", uid).Take(&Message{}).UpdateColumns(
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
	if err := db.Debug().Model(&Message{}).Where("id = ?", uid).Take(&msg).Error; err != nil {
		return &Message{}, err
	}
	return msg, nil
}

func (msg *Message) DeleteMessage(db *gorm.DB, uid uint32) (int64, error) {
	db = db.Debug().Model(&Message{}).Where("id = ?", uid).Take(&Message{}).Delete(&Message{})
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
