package messaging

import (
	"errors"
	"fmt"
	"github.com/sendgrid/sendgrid-go"
	"log"
	"time"
)

var hostURL = "https://api.sendgrid.com"
var sendEmailURL = "/v3/mail/send"

type SendgridService struct {
	key string
}

func (s *SendgridService) ScheduleEmails(emails []EmailPayload, interval int) error {
	sendAt := time.Now()
	request := sendgrid.GetRequest(s.key, sendEmailURL, hostURL)
	request.Method = "POST"
	errCount := 0
	for _, email := range emails {
		sendAt = sendAt.Add(time.Minute * time.Duration(interval))
		request.Body = email.prepareMessagePayload(sendAt.Unix())
		if _, err := sendgrid.API(request); err != nil {
			log.Println(err)
			errCount += 1
		}
	}

	if errCount > 0 {
		return errors.New(fmt.Sprintf("failed to schedule %d emails", errCount))
	}

	return nil
}

func (s *SendgridService) SendMessage(payload EmailPayload) error {
	request := sendgrid.GetRequest(s.key, sendEmailURL, hostURL)
	request.Method = "POST"
	request.Body = payload.prepareMessagePayload()
	if _, err := sendgrid.API(request); err != nil {
		log.Println(err)
		return err
	}
	return nil

}

func (p *EmailPayload) prepareMessagePayload(sendAt ...int64) []byte {
	var sendAtPayload string
	if len(sendAt) > 0 {
		sendAtPayload = fmt.Sprintf(`,"send_at": %d`, sendAt[0])
	}
	bytes := []byte(fmt.Sprintf(` {
		"personalizations": [
			{
				"to": [
					{
						"email": "%s"
					}
				],
				"subject": "%s"
			}
		],
		"from": {
			"email": "email@robertsonlima.com"
		},
		"content": [
			{
				"type": "text/plain",
				"value": "%s"
			}
		]
		%s
	}`, p.To, p.Subject, p.Content, sendAtPayload))
	return bytes
}
