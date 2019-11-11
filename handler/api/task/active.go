package task

import (
	"github.com/fox-one/pkg/logger"
	"github.com/fox-one/pkg/uuid"
	"github.com/gin-gonic/gin"
	"github.com/yiplee/blockquiz/core"
	"github.com/yiplee/blockquiz/handler/api/render"
	"github.com/yiplee/blockquiz/handler/api/request"
	"github.com/yiplee/blockquiz/handler/api/view"
)

func HandleActive(tasks core.TaskStore, courses core.CourseStore, commands core.CommandStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		log := logger.FromContext(ctx)

		task := request.Task(c)
		if !task.IsPending() {
			render.OK(c, view.TaskView(task, nil))
			return
		}

		task.State = core.TaskStateCourse
		if err := tasks.Update(ctx, task); err != nil {
			log.WithError(err).Error("update task")
			render.InternalErrorf(c, "update task failed: %w", err)
			return
		}

		// write show course cmd
		cmd := &core.Command{
			TraceID: uuid.New(),
			UserID:  task.UserID,
			Action:  core.ActionShowCourse,
			Source:  core.CommandSourceOutside,
		}

		if err := commands.Create(ctx, cmd); err != nil {
			log.WithError(err).Error("create show course command")
			render.InternalErrorf(c, "create show course command: %w", err)
			return
		}

		render.OK(c, view.TaskView(task, nil))
	}
}
