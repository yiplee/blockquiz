package render

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/yiplee/blockquiz/handler/api/errors"
)

func ErrorStatus(c *gin.Context, err error, status int) {
	code, msg := errors.Unwrap(err)
	if code == 0 {
		code = status
	}

	view := gin.H{
		"code": code,
		"msg":  msg,
	}

	if gin.IsDebugging() {
		view["hint"] = err.Error()
	}

	c.AbortWithStatusJSON(status, view)
}

func InternalError(c *gin.Context, err error) {
	ErrorStatus(c, err, 500)
}

func InternalErrorf(c *gin.Context, format string, a ...interface{}) {
	ErrorStatus(c, fmt.Errorf(format, a...), 500)
}

func NotFound(c *gin.Context, err error) {
	ErrorStatus(c, err, 404)
}

func NotFoundf(c *gin.Context, format string, a ...interface{}) {
	ErrorStatus(c, fmt.Errorf(format, a...), 404)
}

func Unauthorized(c *gin.Context, err error) {
	ErrorStatus(c, err, 401)
}

func Unauthorizedf(c *gin.Context, format string, a ...interface{}) {
	ErrorStatus(c, fmt.Errorf(format, a...), 401)
}

func Forbidden(c *gin.Context, err error) {
	ErrorStatus(c, err, 403)
}

func Forbiddenf(c *gin.Context, format string, a ...interface{}) {
	ErrorStatus(c, fmt.Errorf(format, a...), 403)
}

func BadRequest(c *gin.Context, err error) {
	ErrorStatus(c, err, 400)
}

func BadRequestf(c *gin.Context, format string, a ...interface{}) {
	ErrorStatus(c, fmt.Errorf(format, a...), 400)
}

func JSON(c *gin.Context, v interface{}, status int) {
	c.AbortWithStatusJSON(status, gin.H{"data": v})
}

func OK(c *gin.Context, v interface{}) {
	JSON(c, v, 200)
}
