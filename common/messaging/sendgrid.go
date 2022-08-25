package messaging

import (
	"fmt"
	"github.com/sendgrid/sendgrid-go"
	"log"
)

var hostURL = "https://api.sendgrid.com"
var sendEmailURL = "/v3/mail/send"

type SendgridService struct {
	key string
}

func (s *SendgridService) SendMessage(to, subject, content string) error {
	request := sendgrid.GetRequest(s.key, sendEmailURL, hostURL)
	request.Method = "POST"
	request.Body = prepareMessagePayload(to, subject, content)
	if _, err := sendgrid.API(request); err != nil {
		log.Println(err)
		return err
	}
	return nil

}

func prepareMessagePayload(to, subject, content string) []byte {
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
	}`, to, subject, content))
	return bytes
}
