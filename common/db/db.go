package db

import (
	"errors"
	"github.com/r0bertson/msg-scheduler/common/models"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"os"
)

type DB interface {
	CreateSession(userID string, userEmail string) (string, error)
	LookupSession(token string) *models.Session

	CreateUser(u *models.User) (*models.User, error)
	UserByEmail(email string) (*models.User, error)
	UserByID(id uint) (*models.User, error)
	Users() (*[]models.User, error)
	DeleteUser(id uint) error

	CreateMessage(message *models.Message) (*models.Message, error)
	UpdateMessage(message *models.Message) (*models.Message, error)
	MessageByID(id uint) (*models.Message, error)
	Messages() (*[]models.Message, error)
	DeleteMessage(id uint) error
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
				log.Fatal(err)
			}
		}
		db, err = gorm.Open(sqlite.Open(url), &gorm.Config{})
	} else {
		db, err = gorm.Open(postgres.Open(url), &gorm.Config{})
	}

	if err != nil {
		log.Fatalln(err)
	}

	db.AutoMigrate(&models.User{}, &models.Message{}, &models.Session{})
	client := &Client{DB: db}
	if result := db.First(&models.Message{}); result.RowsAffected == 0 {
		seedMessages(client)
	}
	return client
}
