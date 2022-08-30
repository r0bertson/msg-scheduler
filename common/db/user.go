package db

import (
	"github.com/r0bertson/msg-scheduler/common/models"
	"github.com/r0bertson/msg-scheduler/common/utils"
	"time"
)

func (c *Client) CreateUser(u *models.User) (*models.User, error) {
	var err error
	//password is still unmasked
	hashedPassword, err := utils.Hash(u.Password)
	if err != nil {
		return nil, err
	}
	u.Password = string(hashedPassword)

	if err = c.DB.Create(&u).Error; err != nil {
		return nil, err
	}
	return u, nil
}

func (c *Client) UserByEmail(email string) (*models.User, error) {
	var user models.User

	if result := c.DB.Where("email = ?", email).First(&user); result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (c *Client) UserByID(id uint) (*models.User, error) {
	var user models.User

	if result := c.DB.First(&user, id); result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (c *Client) Users() (*[]models.User, error) {
	var users []models.User
	result := c.DB.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return &users, nil
}

func (c *Client) InactiveUsers() (*[]models.User, error) {
	return c.UsersByStatus(false)
}

func (c *Client) PendingUsers() (*[]models.User, error) {
	return c.UsersByStatus(true)
}

func (c *Client) UsersByStatus(status bool) (*[]models.User, error) {
	var users []models.User
	result := c.DB.Where("should_send_messages = ?", status).Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}

	return &users, nil
}

func (c *Client) UpdateUserStatus(user *models.User) (*models.User, error) {
	db := c.DB.Model(&models.User{}).Where("id = ?", user.ID).Take(&models.User{}).UpdateColumns(
		map[string]interface{}{
			"should_send_messages": user.ShouldSendMessages,
			"updated_at":           time.Now(),
		},
	)
	if db.Error != nil {
		return nil, db.Error
	}
	// This is the display the updated message
	return c.UserByID(user.ID)
}

func (c *Client) ActivateUsers() error {
	db := c.DB.Debug().Model(&models.User{}).Where("should_send_messages = false").UpdateColumns(
		map[string]interface{}{
			"should_send_messages": true,
			"updated_at":           time.Now(),
		},
	)
	return db.Error
}

func (c *Client) DeleteUser(id uint) error {
	var user models.User
	if result := c.DB.First(&user, id); result.Error != nil {
		return result.Error
	}

	return c.DB.Delete(&user).Error
}
