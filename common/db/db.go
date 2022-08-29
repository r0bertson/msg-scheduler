package db

import (
	"errors"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"msg-scheduler/common/models"
	"os"
)

func Init(url, env string) *gorm.DB {
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
	if result := db.First(&models.Message{}); result.RowsAffected == 0 {
		seedMessages(db)
	}
	return db
}
