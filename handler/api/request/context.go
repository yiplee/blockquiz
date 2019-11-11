package request

import (
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/yiplee/blockquiz/core"
)

const (
	taskKey = "task_context_key"
)

func WithTask(c *gin.Context, task *core.Task) {
	c.Set(taskKey, task)
}

func Task(c *gin.Context) *core.Task {
	return c.MustGet(taskKey).(*core.Task)
}

func Body(c *gin.Context) (body []byte, err error) {
	if cb, ok := c.Get(gin.BodyBytesKey); ok {
		if cbb, ok := cb.([]byte); ok {
			body = cbb
		}
	}

	if body == nil {
		body, err = ioutil.ReadAll(c.Request.Body)
		if err == nil {
			c.Set(gin.BodyBytesKey, body)
		}
	}

	return
}
