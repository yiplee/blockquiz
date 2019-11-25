package task

import (
	"github.com/gin-gonic/gin"
	"github.com/yiplee/blockquiz/core"
	"github.com/yiplee/blockquiz/handler/api/render"
	"github.com/yiplee/blockquiz/handler/api/request"
	"github.com/yiplee/blockquiz/handler/api/view"
)

func HandleFind(courses core.CourseStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		task := request.Task(c)
		render.OK(c, view.TaskView(task, nil))
	}
}
