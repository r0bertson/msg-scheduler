package models

import (
	"gorm.io/gorm"
)

// Session holds the information that associates session cookies with users.
type Session struct {
	gorm.Model
	Token  string
	UserID uint
	Email  string
}

type AuthInfo struct {
	UserID        string
	Authenticated bool
}
