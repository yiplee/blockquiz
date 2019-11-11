package task

import (
	"github.com/fox-one/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/yiplee/blockquiz/core"
	"github.com/yiplee/blockquiz/handler/api/render"
	"github.com/yiplee/blockquiz/handler/api/request"
	"github.com/yiplee/blockquiz/handler/api/view"
)

func HandleFind(courses core.CourseStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		log := logger.FromContext(ctx)

		task := request.Task(c)

		course, err := courses.Find(ctx, task.Course)
		if err != nil {
			log.WithError(err).Error("find course")
			render.InternalErrorf(c, "find course failed: %w", err)
			return
		}

		render.OK(c, view.TaskView(task, course))
	}
}
