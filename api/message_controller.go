package api

import (
	"github.com/gin-gonic/gin"
	"msg-scheduler/common/models"
)

func (h handler) GetMessage(c *gin.Context) (interface{}, error) {
	var msg models.Message

	if result := h.DB.First(&msg, c.Param("id")); result.Error != nil {
		return NotFoundWithMessage(c, result.Error.Error())
	}

	return msg, nil
}

func (h handler) GetMessages(c *gin.Context) (interface{}, error) {
	var msgs []models.Message
	if result := h.DB.Find(&msgs); result.Error != nil {
		return NotFoundWithMessage(c, result.Error.Error())
	}

	return msgs, nil
}

func (h handler) DeleteMessage(c *gin.Context) (interface{}, error) {
	var msg models.Message
	if result := h.DB.First(&msg, c.Param("id")); result.Error != nil {
		return NotFoundWithMessage(c, result.Error.Error())
	}

	h.DB.Delete(&msg)
	return NoContent(c)
}
func (h handler) CreateMessage(c *gin.Context) (interface{}, error) {
	// getting request's body
	body := MessageOperationsRequestBody{}
	if err := c.BindJSON(&body); err != nil {
		return BadRequest(c, err.Error())
	}
	//validate payload
	if err := body.Validate(models.Create); err != nil {
		return BadRequest(c, err.Error())
	}

	var message models.Message

	message.Content = body.Content
	message.Subject = body.Subject

	savedUser, err := message.SaveMessage(h.DB)
	if err != nil {
		return nil, err
	}

	return savedUser, err
}

func (h handler) UpdateMessage(c *gin.Context) (interface{}, error) {
	var message models.Message
	if result := h.DB.First(&message, c.Param("id")); result.Error != nil {
		return NotFoundWithMessage(c, result.Error.Error())
	}

	// getting request's body
	body := MessageOperationsRequestBody{}
	if err := c.BindJSON(&body); err != nil {
		return BadRequest(c, err.Error())
	}

	//validate payload
	if err := body.Validate(models.Create); err != nil {
		return BadRequest(c, err.Error())
	}

	message.Subject = body.Subject
	message.Content = body.Content

	h.DB.Save(&message)
	return message, nil
}
