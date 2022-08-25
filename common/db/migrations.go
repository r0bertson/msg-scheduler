package db

import (
	"fmt"
	"gorm.io/gorm"
	"log"
	"math/rand"
	"msg-scheduler/common/models"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func getDefaultMessages() []models.Message {
	msgs := []models.Message{}
	for i := 1; i <= 10; i++ {
		msgs = append(msgs, models.Message{Subject: fmt.Sprintf("Test message number %d", i), Content: RandStringRunes(50)})
	}
	return msgs
}

func seedMessages(db *gorm.DB) {
	for _, msg := range getDefaultMessages() {
		if _, err := msg.SaveMessage(db); err != nil {
			log.Fatal(err)
		}
	}
}
