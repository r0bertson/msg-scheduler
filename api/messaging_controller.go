package api

import (
	"github.com/gin-gonic/gin"
	"msg-scheduler/common/messaging"
)

func (h handler) SendTestMessage(c *gin.Context) (interface{}, error) {
	// getting request's body
	msg := messaging.EmailPayload{}
	if err := c.BindJSON(&msg); err != nil {
		return BadRequest(c, err.Error())
	}
	return nil, h.msgService.SendMessage(msg)
}
