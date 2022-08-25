package messaging

import (
	"log"
	"strings"
)

type MsgService interface {
	SendMessage(to, subject, content string) error
}

func Init(service, key string) MsgService {
	switch strings.ToLower(service) {
	case "sendgrid":
		sendgrid := SendgridService{key: key}
		return &sendgrid
	default:
		log.Fatalln("messaging service not implemented")
	}
	return nil
}
