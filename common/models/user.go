package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email              string `gorm:"size:100;not null;unique"`
	Password           string `gorm:"size:100;not null;" json:"-"`
	ShouldSendMessages bool
}

type Stats struct {
	SentMessages []uint
}

type UserStats struct {
	User
	Stats Stats
}
