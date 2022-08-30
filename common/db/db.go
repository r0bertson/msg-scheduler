package db

import (
	"errors"
	"github.com/r0bertson/msg-scheduler/common/models"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
)

type DB interface {
	CreateSession(userID string, userEmail string) (string, error)
	LookupSession(token string) *models.Session

	CreateUser(u *models.User) (*models.User, error)
	UpdateUserStatus(user *models.User) (*models.User, error)
	UserByEmail(email string) (*models.User, error)
	UserByID(id uint) (*models.User, error)
	Users() (*[]models.User, error)
	PendingUsers() (*[]models.User, error)
	DeleteUser(id uint) error

	CreateMessage(message *models.Message) (*models.Message, error)
	UpdateMessage(message *models.Message) (*models.Message, error)
	MessageByID(id uint) (*models.Message, error)
	Messages() (*[]models.Message, error)
	DeleteMessage(id uint) error

	CreateSentMessage(sent *models.SentMessage) (*models.SentMessage, error)
	MessagesSentFor(ids []uint) (*[]models.SentMessage, error)
	MessagesSent() (*[]models.SentMessage, error)
	MessagesByUserID(id uint) (*[]models.SentMessage, error)
}

type Client struct {
	DB *gorm.DB
}

func Init(url, env string) *Client {
	var db *gorm.DB
	var err error

	if env == "local" {
		_, err := os.Stat(url)
		if errors.Is(err, os.ErrNotExist) {
			if _, err := os.Create(url); err != nil {
				log.Fatal().Msg(err.Error())
			}
		}
		db, err = gorm.Open(sqlite.Open(url), &gorm.Config{})
	} else {
		db, err = gorm.Open(postgres.Open(url), &gorm.Config{})
	}

	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	db.AutoMigrate(&models.User{}, &models.Message{}, &models.Session{}, &models.SentMessage{})
	client := &Client{DB: db}
	if result := db.First(&models.Message{}); result.RowsAffected == 0 {
		seedMessages(client)
	}
	return client
}
