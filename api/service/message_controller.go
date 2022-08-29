package service

import (
	"github.com/gin-gonic/gin"
	"github.com/r0bertson/msg-scheduler/common/models"
	"github.com/r0bertson/msg-scheduler/common/utils"
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
	msgID, err := utils.UintID(c.Param("id"))
	if err != nil {
		return BadRequest(c, err.Error())
	}

	msg, err := h.DB.MessageByID(msgID)
	if err != nil {
		return NotFoundWithMessage(c, err.Error())
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
// @Failure      500  {object}  ErrResp
// @Router       /messages [get]
func (h handler) GetMessages(c *gin.Context) (interface{}, error) {
	return h.DB.Messages()
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
	msgID, err := utils.UintID(c.Param("id"))
	if err != nil {
		return BadRequest(c, err.Error())
	}

	if err = h.DB.DeleteMessage(msgID); err != nil {
		return NotFoundWithMessage(c, err.Error())
	}

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

	message := &models.Message{
		Content: body.Content,
		Subject: body.Subject,
	}

	message, err := h.DB.CreateMessage(message)
	if err != nil {
		return nil, err
	}

	return message, err
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
	msgID, err := utils.UintID(c.Param("id"))
	if err != nil {
		return BadRequest(c, err.Error())
	}

	msg, err := h.DB.MessageByID(msgID)
	if err != nil {
		return NotFoundWithMessage(c, "message not found")
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

	msg.Subject = body.Subject
	msg.Content = body.Content

	return h.DB.UpdateMessage(msg)
}
