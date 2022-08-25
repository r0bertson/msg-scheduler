package api

import "github.com/gin-gonic/gin"

func (h handler) SendTestMessage(c *gin.Context) (interface{}, error) {
	// getting request's body
	msg := TestMessageOperationsRequestBody{}
	if err := c.BindJSON(&msg); err != nil {
		return BadRequest(c, err.Error())
	}
	return nil, h.msgService.SendMessage(msg.To, msg.Subject, msg.Content)
}
