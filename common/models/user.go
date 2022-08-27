package models

import (
	"gorm.io/gorm"
	"msg-scheduler/common/utils"
)

type User struct {
	gorm.Model
	Email    string `gorm:"size:100;not null;unique" json:"Email"`
	Password string `gorm:"size:100;not null;" json:"-"`
}

func (u *User) hashPassword() error {
	hashedPassword, err := utils.Hash(u.Password)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) SaveUser(db *gorm.DB) (*User, error) {
	var err error
	u.hashPassword()
	err = db.Debug().Create(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}
