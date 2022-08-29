package db

import (
	"github.com/r0bertson/msg-scheduler/common/models"
	"github.com/r0bertson/msg-scheduler/common/utils"
)

func (c *Client) CreateUser(u *models.User) (*models.User, error) {
	var err error
	//password is still unmasked
	hashedPassword, err := utils.Hash(u.Password)
	if err != nil {
		return nil, err
	}
	u.Password = string(hashedPassword)

	if err = c.DB.Debug().Create(&u).Error; err != nil {
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
	if result := c.DB.Find(&users); result.Error != nil {
		return nil, result.Error
	}
	return &users, nil
}

func (c *Client) DeleteUser(id uint) error {
	var user models.User
	if result := c.DB.First(&user, id); result.Error != nil {
		return result.Error
	}

	return c.DB.Delete(&user).Error
}
