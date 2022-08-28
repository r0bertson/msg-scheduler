package api

import (
	"github.com/gin-gonic/gin"
	"msg-scheduler/common/models"
)

// GetMessage godoc
// @Summary      fetches a specific message
// @Description  fetches a specific message
// @Tags         messages
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "MessageID"
// @Success      200  {object}  ErrResp
// @Failure      404  {object}  ErrResp
// @Router       /messages/{id} [get]
func (h handler) GetMessage(c *gin.Context) (interface{}, error) {
	var msg models.Message

	if result := h.DB.First(&msg, c.Param("id")); result.Error != nil {
		return NotFoundWithMessage(c, result.Error.Error())
	}

	return msg, nil
}

// GetMessages godoc
// @Summary      fetches all messages
// @Description  fetches all messages
// @Tags         messages
// @Accept       json
// @Produce      json
// @Success      200  {array}  ErrResp
// @Failure      404  {object}  ErrResp
// @Router       /messages [get]
func (h handler) GetMessages(c *gin.Context) (interface{}, error) {
	var msgs []models.Message
	if result := h.DB.Find(&msgs); result.Error != nil {
		return NotFoundWithMessage(c, result.Error.Error())
	}

	return msgs, nil
}

// DeleteMessage godoc
// @Summary      delete a message
// @Description  delete a message
// @Tags         messages
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "MessageID"
// @Success      204
// @Failure      404  {object}  ErrResp
// @Router       /messages/{id} [delete]
func (h handler) DeleteMessage(c *gin.Context) (interface{}, error) {
	var msg models.Message
	if result := h.DB.First(&msg, c.Param("id")); result.Error != nil {
		return NotFoundWithMessage(c, result.Error.Error())
	}

	h.DB.Delete(&msg)
	return NoContent(c)
}

// CreateMessage godoc
// @Summary      Creates a new message
// @Description  creates a new message
// @Tags         messages
// @Accept       json
// @Produce      json
// @Param  data body MessageOperationsRequestBody true "message operations request"
// @Success      200  {object}  models.User
// @Failure      400  {object}  ErrResp
// @Router       /messages [post]
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

// UpdateMessage godoc
// @Summary      updates a message
// @Description  updates a message
// @Tags         messages
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "MessageId"
// @Param  		 data body MessageOperationsRequestBody true "message operations request"
// @Success      200  {object}  models.User
// @Failure      400  {object}  ErrResp
// @Router       /messages/{id} [post]
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
