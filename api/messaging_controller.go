package api

import (
	"github.com/gin-gonic/gin"
	"msg-scheduler/common/messaging"
)

func (h handler) SendTestMessage(c *gin.Context) (interface{}, error) {
	// getting request's body
	msg := SendEmailRequestBody{}
	if err := c.BindJSON(&msg); err != nil {
		return BadRequest(c, err.Error())
	}
	if msg.Timeout != nil {
		return nil, h.msgService.ScheduleEmails([]messaging.EmailPayload{msg.Payload}, *msg.Timeout)
	}
	return nil, h.msgService.SendMessage(msg.Payload)
}
