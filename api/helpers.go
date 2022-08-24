package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func BasicHandler(inner func(c *gin.Context) (interface{}, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Run internal handler: returns a parseable result and an error, either of which may be nil.
		result, err := inner(c)

		// Propagate Go errors as "500 Internal Server Error" responses.
		if err != nil {
			log.Printf("handling %q: %v", c.Request.RequestURI, err)
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		// No response body, so internal handler dealt with response setup.
		if result == nil {
			return
		}
		c.JSON(http.StatusOK, result)
	}
}

// ErrResp JSON error response.
type ErrResp struct {
	Message string `json:"message"`
}

// BadRequest sets up an HTTP 400 Bad Request with a given error message and returns
// a (nil, nil) pair used by BasicHandler to signal that the response has been dealt with.
func BadRequest(c *gin.Context, msg string) (interface{}, error) {
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusBadRequest, ErrResp{msg})
	return nil, nil
}

// NotFound sets up an HTTP 404 Not Found and returns the (nil, nil) pair used by BasicHandler
// to signal that the response has been dealt with.
func NotFound(c *gin.Context) (interface{}, error) {
	c.AbortWithStatus(http.StatusNotFound)
	return nil, nil
}

// NotFoundWithMessage sets up an HTTP 404 Not Found with a given error message and returns the
// (nil, nil) pair used by BasicHandler to signal that the response has been dealt with.
func NotFoundWithMessage(c *gin.Context, msg string) (interface{}, error) {
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusNotFound, ErrResp{msg})
	return nil, nil
}

// Forbidden sets up an HTTP 403 Forbidden and returns the (nil, nil) pair used by BasicHandler
// to signal that the response has been dealt with.
func Forbidden(c *gin.Context) (interface{}, error) {
	c.AbortWithStatus(http.StatusForbidden)
	return nil, nil
}

// Unauthorized sets up an HTTP 401 StatusUnauthorized and returns the (nil, nil)
// pair used by BasicHandler to signal that the response has been dealt with.
func Unauthorized(c *gin.Context, msg string) (interface{}, error) {
	rsp := ErrResp{msg}
	body, _ := json.Marshal(rsp)
	c.Header("Content-Type", "application/json")
	c.AbortWithStatusJSON(http.StatusUnauthorized, body)
	return nil, nil
}

// NoContent sets up an HTTP 204 No Content and returns the (nil, nil) pair used by
// BasicHandler to signal that the response has been dealt with.
func NoContent(c *gin.Context) (interface{}, error) {
	c.Status(http.StatusNoContent)
	return nil, nil
}
