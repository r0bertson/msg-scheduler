package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string `gorm:"size:100;not null;unique"`
	Password string `gorm:"size:100;not null;" json:"-"`
}
