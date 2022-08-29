package db

import (
	"github.com/r0bertson/msg-scheduler/common/models"
)

func (c *Client) CreateSentMessage(sent *models.SentMessage) (*models.SentMessage, error) {
	if err := c.DB.Create(sent).Error; err != nil {
		return nil, err
	}
	return sent, nil
}
func (c *Client) MessagesByUserID(id uint) (*[]models.SentMessage, error) {
	var messages []models.SentMessage
	if result := c.DB.Where("user_id = ?", id).Find(&messages); result.Error != nil {
		return nil, result.Error
	}
	return &messages, nil
}

func (c *Client) MessagesSentFor(ids []uint) (*[]models.SentMessage, error) {
	var messagesSent []models.SentMessage
	if err := c.DB.Raw("SELECT * from sent_messages WHERE user_id IN (?)", ids).Scan(&messagesSent).Error; err != nil {
		return nil, err
	}
	return &messagesSent, nil
}
