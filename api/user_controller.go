package api

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"msg-scheduler/common/messaging"
	"msg-scheduler/common/models"
)

type handler struct {
	DB         *gorm.DB
	msgService messaging.MsgService
}

func (h handler) GetUser(c *gin.Context) (interface{}, error) {
	var user models.User
	if result := h.DB.First(&user, c.Param("id")); result.Error != nil {
		return NotFoundWithMessage(c, result.Error.Error())
	}

	return user, nil
}

func (h handler) GetUsers(c *gin.Context) (interface{}, error) {
	var users []models.User
	if result := h.DB.Find(&users); result.Error != nil {
		return NotFoundWithMessage(c, result.Error.Error())
	}

	return users, nil
}

func (h handler) DeleteUser(c *gin.Context) (interface{}, error) {
	var user models.User
	if result := h.DB.First(&user, c.Param("id")); result.Error != nil {
		return NotFoundWithMessage(c, result.Error.Error())
	}

	h.DB.Delete(&user)
	return NoContent(c)
}

func (h handler) CreateUser(c *gin.Context) (interface{}, error) {
	// getting request's body
	body := UserOperationsRequestBody{}
	if err := c.BindJSON(&body); err != nil {
		return BadRequest(c, err.Error())
	}
	//validate payload
	if err := body.Validate(models.Create); err != nil {
		return BadRequest(c, err.Error())
	}

	var user models.User
	if result := h.DB.Where("email = ?", body.Email).First(&user); result.Error == nil {
		return BadRequest(c, "user already created")
	}

	user.Email = body.Email
	user.Password = body.Password

	savedUser, err := user.SaveUser(h.DB)
	if err != nil {
		return nil, err
	}
	go func() {
		var msgs []models.Message
		if result := h.DB.Limit(10).Find(&msgs); result.Error != nil {
			return //this can "safely" fail because a routine will pick this up later
		}
		var emails []messaging.EmailPayload
		for _, msg := range msgs {
			emails = append(emails, messaging.EmailPayload{To: user.Email, Subject: msg.Subject, Content: msg.Content})
		}
		h.msgService.ScheduleEmails(emails, 1)
	}()
	return savedUser, err
}

func (h handler) UpdateUser(c *gin.Context) (interface{}, error) {
	var user models.User
	if result := h.DB.First(&user, c.Param("id")); result.Error != nil {
		return NotFoundWithMessage(c, result.Error.Error())
	}

	// getting request's body
	body := UserOperationsRequestBody{}
	if err := c.BindJSON(&body); err != nil {
		return BadRequest(c, err.Error())
	}

	//validate payload
	if err := body.Validate(models.Create); err != nil {
		return BadRequest(c, err.Error())
	}

	user.Password = body.Password
	user.Email = body.Email

	user.UpdateUser(h.DB)
	return user, nil
}
