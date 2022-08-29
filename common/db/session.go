package db

import (
	"github.com/r0bertson/msg-scheduler/common/models"
	"github.com/r0bertson/msg-scheduler/common/utils"
)

// CreateSession generates a new session token and stores it along with the associated user information.
func (c *Client) CreateSession(userID uint, userEmail string) (*models.Session, error) {
	session := models.Session{
		Token:  utils.NewBareID(24),
		UserID: userID,
		Email:  userEmail,
	}

	if err := c.DB.Debug().Create(&session).Error; err != nil {
		return nil, err
	}
	return &session, nil
}

// LookupSession checks a session token and returns user info if the session is known.
func (c *Client) LookupSession(token string) *models.Session {
	session := models.Session{}
	if result := c.DB.Where("token = ?", token).First(&session); result.Error != nil {
		return nil
	}
	return &session
}
