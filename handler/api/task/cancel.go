package task

import (
	"github.com/fox-one/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/yiplee/blockquiz/core"
	"github.com/yiplee/blockquiz/handler/api/render"
	"github.com/yiplee/blockquiz/handler/api/request"
	"github.com/yiplee/blockquiz/handler/api/view"
)

func HandleCancel(tasks core.TaskStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		log := logger.FromContext(ctx)

		task := request.Task(c)
		if task.IsDone() {
			render.OK(c, view.TaskView(task, nil))
			return
		}

		task.State = core.TaskStateCancelled
		if err := tasks.Update(ctx, task); err != nil {
			log.WithError(err).Error("update task")
			render.InternalErrorf(c, "update task failed: %w", err)
			return
		}

		render.OK(c, view.TaskView(task, nil))
	}
}
