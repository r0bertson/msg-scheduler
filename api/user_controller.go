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

	h.DB.Save(&user)
	return user, nil
}
