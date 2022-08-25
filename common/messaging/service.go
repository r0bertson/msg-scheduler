package messaging

import (
	"log"
	"strings"
)

type EmailPayload struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Content string `json:"content"`
}

type MsgService interface {
	SendMessage(payload EmailPayload) error
	ScheduleEmails(emails []EmailPayload, interval int) error
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
