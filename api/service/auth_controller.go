package service

import (
	"github.com/gin-gonic/gin"
	"github.com/r0bertson/msg-scheduler/common/models"
	"github.com/r0bertson/msg-scheduler/common/utils"
)

// Login godoc
// @Summary      Authenticates a user
// @Description  Authenticates a user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param  data body UserOperationsRequestBody true "email and pwd struct"
// @Success      200  {object}  models.Session
// @Failure      400  {object}  ErrResp
// @Failure      401  {object}  ErrResp
// @Failure      404  {object}  ErrResp
// @Router       /auth/login [post]
func (h handler) Login(c *gin.Context) (interface{}, error) {
	// Decode request body to get email and password.
	var req UserOperationsRequestBody
	if err := c.ShouldBindJSON(&req); err != nil {
		return BadRequest(c, err.Error())
	}
	user, err := h.DB.UserByEmail(req.Email)
	if err != nil {
		return NotFoundWithMessage(c, err.Error())
	}

	if !utils.CompareHashedKeys(user.Password, req.Password) {
		return Unauthorized(c, "invalid password or email")
	}

	// Create a session for the user.
	session, err := h.DB.CreateSession(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	// Set the session cookie and return the user information as a JSON response.
	c.SetCookie("session", session.Token, 216000, "/", "", false, true)

	return session, nil

}

func (h handler) getAuthentication(c *gin.Context) *models.Session {
	if token, err := c.Cookie("session"); err == nil {
		return h.DB.LookupSession(token)
	}
	return nil
}

func (h handler) userHasPermission(c *gin.Context, userId uint) bool {
	auth := h.getAuthentication(c)
	if auth == nil {
		NotFound(c)
		return false
	}
	if auth.ID != userId {
		Unauthorized(c, "authenticated user has no access to this resource")
		return false
	}
	return true
}
