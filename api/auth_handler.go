package api

import (
	"github.com/gin-gonic/gin"
	"github.com/msg-scheduler/common/models"
	"github.com/msg-scheduler/common/utils"
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
	var user models.User
	if result := h.DB.First(&user, c.Param("id")); result.Error != nil {
		return NotFoundWithMessage(c, result.Error.Error())
	}

	if !utils.CompareHashedKeys(user.Password, req.Password) {
		return Unauthorized(c, "invalid password or email")
	}

	// Create a session for the user.
	session, err := user.CreateSession(h.DB)
	if err != nil {
		return nil, err
	}

	// Set the session cookie and return the user information as a JSON response.
	c.SetCookie("session", session.Token, 216000, "/", "", false, true)

	return session, nil

}

func (h handler) getAuthentication(c *gin.Context) *models.Session {
	session := models.Session{}
	if token, err := c.Cookie("session"); err == nil {
		return session.LookupSession(h.DB, token)
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
