package api

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"msg-scheduler/common/messaging"
	"msg-scheduler/common/models"
	"msg-scheduler/common/utils"
)

type handler struct {
	DB         *gorm.DB
	msgService messaging.MsgService
}

// GetUser godoc
// @Summary      Get a specific user
// @Description  get string by ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "UserID"
// @Success      200  {object}  models.User
// @Failure      400  {object}  ErrResp
// @Failure      401  {object}  ErrResp
// @Failure      404  {object}  ErrResp
// @Router       /users/{id} [get]
func (h handler) GetUser(c *gin.Context) (interface{}, error) {
	userId, err := utils.UintID(c.Param("id"))
	if err != nil {
		return BadRequest(c, err.Error())
	}

	if !h.userHasPermission(c, userId) {
		return nil, nil //already handled
	}

	var user models.User
	if result := h.DB.First(&user, c.Param("id")); result.Error != nil {
		return NotFoundWithMessage(c, result.Error.Error())
	}

	return user, nil
}

// GetMe godoc
// @Summary      Get the authenticated user
// @Description  get string by ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.User
// @Failure      404  {object}  ErrResp
// @Router       /users/me [get]
func (h handler) GetMe(c *gin.Context) (interface{}, error) {
	auth := h.getAuthentication(c)
	if auth == nil {
		return NotFound(c)
	}

	users := models.User{Model: gorm.Model{ID: auth.UserID}}
	if result := h.DB.Find(&users); result.Error != nil {
		return NotFoundWithMessage(c, result.Error.Error())
	}

	return users, nil
}

// GetUsers godoc
// @Summary      Get all users
// @Description  get all users
// @Tags         users
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.User
// @Failure      404  {object}  ErrResp
// @Router       /users [get]
func (h handler) GetUsers(c *gin.Context) (interface{}, error) {
	var users []models.User
	if result := h.DB.Find(&users); result.Error != nil {
		return NotFoundWithMessage(c, result.Error.Error())
	}

	return users, nil
}

// DeleteUser godoc
// @Summary      Deletes a specific user
// @Description  get string by ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "UserID"
// @Success      204
// @Failure      400  {object}  ErrResp
// @Failure      401  {object}  ErrResp
// @Failure      404  {object}  ErrResp
// @Failure      500  {object}  ErrResp
// @Router       /users/{id} [delete]
func (h handler) DeleteUser(c *gin.Context) (interface{}, error) {
	userId, err := utils.UintID(c.Param("id"))
	if err != nil {
		return BadRequest(c, err.Error())
	}

	if !h.userHasPermission(c, userId) {
		return nil, nil //already handled
	}
	var user models.User
	if result := h.DB.First(&user); result.Error != nil {
		return NotFoundWithMessage(c, result.Error.Error())
	}

	h.DB.Delete(&user)
	return NoContent(c)
}

// CreateUser godoc
// @Summary      Creates a new user
// @Description  creates a new user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param  data body UserOperationsRequestBody true "email and pwd struct"
// @Success      200  {object}  models.User
// @Failure      400  {object}  ErrResp
// @Failure      401  {object}  ErrResp
// @Failure      404  {object}  ErrResp
// @Failure      500  {object}  ErrResp
// @Router       /users [post]
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
