package models

import (
	"gorm.io/gorm"
	"msg-scheduler/common/utils"
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

// CreateSession generates a new session token and stores it along with the associated user information.
func (u *User) CreateSession(db *gorm.DB) (*Session, error) {
	session := Session{
		Token:  utils.NewBareID(24),
		UserID: u.ID,
		Email:  u.Email,
	}

	if err := db.Debug().Create(&session).Error; err != nil {
		return nil, err
	}

	return &session, nil
}

// LookupSession checks a session token and returns user info if the session is known.
func (s *Session) LookupSession(db *gorm.DB, token string) *Session {
	session := Session{}
	if result := db.Where("token = ?", token).First(&session); result.Error != nil {
		return nil
	}
	return &session
}
