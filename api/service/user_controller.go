package service

import (
	"github.com/gin-gonic/gin"
	"github.com/r0bertson/msg-scheduler/common/db"
	"github.com/r0bertson/msg-scheduler/common/messaging"
	"github.com/r0bertson/msg-scheduler/common/models"
	"github.com/r0bertson/msg-scheduler/common/utils"
)

type handler struct {
	DB         *db.Client
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

	user, err := h.DB.UserByID(userId)
	if err != nil {
		return NotFoundWithMessage(c, err.Error())
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
	return h.DB.UserByID(auth.UserID)
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
	return h.DB.Users()
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

	if err = h.DB.DeleteUser(userId); err != nil {
		return nil, err
	}

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

	user, err := h.DB.UserByEmail(body.Email)
	if err == nil {
		return BadRequest(c, "user already created")
	}
	user = &models.User{Email: body.Email, Password: body.Password, ShouldSendMessages: true}

	user, err = h.DB.CreateUser(user)
	if err != nil {
		return nil, err
	}

	return user, err
}
