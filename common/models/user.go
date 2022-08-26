package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"time"
)

type User struct {
	gorm.Model
	Email    string `gorm:"size:100;not null;unique" json:"Email"`
	Password string `gorm:"size:100;not null;" json:"-"`
}

func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func (u *User) hashPassword() error {
	hashedPassword, err := Hash(u.Password)
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

func (u *User) UpdateUser(db *gorm.DB) (*User, error) {

	// hash the password
	if err := u.hashPassword(); err != nil {
		log.Fatal(err)
	}
	db = db.Debug().Model(&User{}).Where("id = ?", u.ID).Take(&User{}).UpdateColumns(
		map[string]interface{}{
			"password":   u.Password,
			"email":      u.Email,
			"updated_at": time.Now(),
		},
	)
	if db.Error != nil {
		return &User{}, db.Error
	}
	// This is the display the updated user
	if err := db.Debug().Model(&User{}).Where("id = ?", u.ID).Take(&u).Error; err != nil {
		return &User{}, err
	}
	return u, nil
}
