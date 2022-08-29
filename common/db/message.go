package db

import (
	"github.com/r0bertson/msg-scheduler/common/models"
	"time"
)

func (c *Client) CreateMessage(message *models.Message) (*models.Message, error) {
	if err := c.DB.Debug().Create(message).Error; err != nil {
		return nil, err
	}
	return message, nil
}

func (c *Client) UpdateMessage(message *models.Message) (*models.Message, error) {
	db := c.DB.Debug().Model(&models.Message{}).Where("id = ?", message.ID).Take(&models.Message{}).UpdateColumns(
		map[string]interface{}{
			"content":    message.Content,
			"subject":    message.Subject,
			"updated_at": time.Now(),
		},
	)
	if db.Error != nil {
		return nil, db.Error
	}
	// This is the display the updated message
	return c.MessageByID(message.ID)
}

func (c *Client) MessageByID(id uint) (*models.Message, error) {
	var msg models.Message

	if result := c.DB.First(&msg, id); result.Error != nil {
		return nil, result.Error
	}
	return &msg, nil
}

func (c *Client) Messages() (*[]models.Message, error) {
	var messages []models.Message
	if result := c.DB.Find(&messages); result.Error != nil {
		return nil, result.Error
	}
	return &messages, nil
}

func (c *Client) DeleteMessage(id uint) error {
	var msg models.Message
	if result := c.DB.First(&msg, id); result.Error != nil {
		return result.Error
	}

	return c.DB.Delete(&msg).Error
}
