package service

import (
	"github.com/gin-gonic/gin"
	"github.com/r0bertson/msg-scheduler/common/messaging"
)

// SendTestMessage godoc
// @Summary      Sends a test message to an email
// @Description  Sends a test message to an email
// @Tags         messaging
// @Accept       json
// @Produce      json
// @Param  data body SendEmailRequestBody true "send test email request"
// @Success      200
// @Failure      400  {object}  ErrResp
// @Router       /messaging/send [post]
func (h handler) SendTestMessage(c *gin.Context) (interface{}, error) {
	// getting request's body
	msg := SendEmailRequestBody{}
	if err := c.BindJSON(&msg); err != nil {
		return BadRequest(c, err.Error())
	}
	if msg.Timeout != nil && *msg.Timeout > 0 {
		return nil, h.msgService.ScheduleEmails([]messaging.EmailPayload{msg.Payload}, *msg.Timeout)
	}
	return nil, h.msgService.SendMessage(msg.Payload)
}
