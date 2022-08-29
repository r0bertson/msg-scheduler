package db

import (
	"fmt"
	"github.com/r0bertson/msg-scheduler/common/models"
	"log"
	"math/rand"
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

func seedMessages(db *Client) {
	for _, msg := range getDefaultMessages() {
		if _, err := db.CreateMessage(&msg); err != nil {
			log.Fatal(err)
		}
	}
}
